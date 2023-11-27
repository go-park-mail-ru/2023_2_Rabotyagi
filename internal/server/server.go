package server

import (
	"context"
	"fmt"
	categoryrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/repository"
	categoryusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/usecases"
	cityrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/repository"
	cityusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery/mux"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/repository"
	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/repository"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"google.golang.org/grpc"
	"net/http"
	"strings"
	"time"
)

const (
	pathCertFile = "/etc/ssl/goods-galaxy.ru.crt"
	pathKeyFile  = "/etc/ssl/goods-galaxy.ru.key"

	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

//nolint:funlen
func (s *Server) Run(config *config.Config) error {
	baseCtx := context.Background()

	grcpConnAuth, err := grpc.Dial(
		":8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println(err)

		return err
	}
	defer grcpConnAuth.Close()

	authGrpcService := auth.NewSessionMangerClient(grcpConnAuth)

	pool, err := repository.NewPgxPool(baseCtx, config.URLDataBase)
	if err != nil {
		return err //nolint:wrapcheck
	}

	logger, err := my_logger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer logger.Sync()

	productStorage, err := productrepo.NewProductStorage(pool)
	if err != nil {
		return err
	}

	basketService, err := usecases.NewBasketService(productStorage)
	if err != nil {
		return err
	}

	favouriteService, err := usecases.NewFavouriteService(productStorage)
	if err != nil {
		return err
	}

	productService, err := usecases.NewProductService(productStorage, *basketService, *favouriteService)
	if err != nil {
		return err
	}

	userStorage, err := userrepo.NewUserStorage(pool)
	if err != nil {
		return err
	}

	userService, err := userusecases.NewUserService(userStorage)
	if err != nil {
		return err
	}

	categoryStorage, err := categoryrepo.NewCategoryStorage(pool)
	if err != nil {
		return err
	}

	categoryService, err := categoryusecases.NewCategoryService(categoryStorage)
	if err != nil {
		return err
	}

	cityStorage, err := cityrepo.NewCityStorage(pool)
	if err != nil {
		return err
	}

	cityService, err := cityusecases.NewCityService(cityStorage)
	if err != nil {
		return err
	}

	handler, err := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin,
		config.Schema, config.PortServer),
		userService, productService, categoryService, cityService, authGrpcService, logger)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.PortServer,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    basicTimeout,
		WriteTimeout:   basicTimeout,
	}

	logger.Infof("Start server:%s", config.PortServer)

	if config.ProductionMode {
		return s.httpServer.ListenAndServeTLS(pathCertFile, pathKeyFile)
	}

	chCloseRefreshing := make(chan struct{})

	// don`t want use chCloseRefreshing secret now
	jwt.StartRefreshingSecret(time.Hour*jwt.TimeRefreshSecretInHours, chCloseRefreshing)

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

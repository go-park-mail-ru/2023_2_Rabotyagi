package server

import (
	"context"
	categoryusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/usecases"
	cityusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/city/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
	fileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"google.golang.org/grpc"
	"net/http"
	"strings"
	"time"

	categoryrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/repository"
	cityrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/city/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/mux"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/repository"
	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
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

	grpcConnFileService, err := grpc.Dial(
		config.AddressFileServiceGrpc,
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	defer grpcConnFileService.Close()

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

	basketService, err := productusecases.NewBasketService(productStorage)
	if err != nil {
		return err
	}

	favouriteService, err := productusecases.NewFavouriteService(productStorage)
	if err != nil {
		return err
	}

	fileServiceClient := fileservice.NewFileServiceClient(grpcConnFileService)

	productService, err := productusecases.NewProductService(productStorage,
		*basketService, *favouriteService, fileServiceClient)
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
		userService, productService, categoryService, cityService, logger)
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

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	categoryrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/repository"
	categoryusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/usecases"
	cityrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/repository"
	cityusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/config"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery/mux"
	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/repository"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	fileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

//nolint:funlen
func (s *Server) Run(config *config.Config) error { //nolint:cyclop
	baseCtx := context.Background()

	grcpConnAuth, err := grpc.Dial(
		config.AddressAuthServiceGrpc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println(err)

		return err //nolint:wrapcheck
	}
	defer grcpConnAuth.Close()

	grpcConnFileService, err := grpc.Dial(
		config.AddressFileServiceGrpc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err //nolint:wrapcheck
	}
	defer grpcConnFileService.Close()

	authGrpcService := auth.NewSessionMangerClient(grcpConnAuth)

	fileServiceClient := fileservice.NewFileServiceClient(grpcConnFileService)

	pool, err := repository.NewPgxPool(baseCtx, config.URLDataBase)
	if err != nil {
		return err //nolint:wrapcheck
	}

	logger, err := mylogger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer logger.Sync() //nolint:errcheck

	productStorage, err := productrepo.NewProductStorage(pool)
	if err != nil {
		return err //nolint:wrapcheck
	}

	basketService, err := usecases.NewBasketService(productStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	favouriteService, err := usecases.NewFavouriteService(productStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	premiumService, err := usecases.NewPremiumService(productStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	commentService, err := usecases.NewCommentService(productStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	productService, err := usecases.NewProductService(productStorage, basketService, favouriteService,
		premiumService, commentService, fileServiceClient)
	if err != nil {
		return err //nolint:wrapcheck
	}

	userStorage, err := userrepo.NewUserStorage(pool)
	if err != nil {
		return err //nolint:wrapcheck
	}

	userService, err := userusecases.NewUserService(userStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	categoryStorage, err := categoryrepo.NewCategoryStorage(pool)
	if err != nil {
		return err //nolint:wrapcheck
	}

	categoryService, err := categoryusecases.NewCategoryService(categoryStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	cityStorage, err := cityrepo.NewCityStorage(pool)
	if err != nil {
		return err //nolint:wrapcheck
	}

	cityService, err := cityusecases.NewCityService(cityStorage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	handler, err := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin,
		config.Schema, config.PortServer, config.MainServiceName,
		config.PremiumShopID, config.PremiumShopSecret, config.PathCertFile),
		userService, productService, categoryService, cityService, authGrpcService, logger)
	if err != nil {
		return err //nolint:wrapcheck
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
		return s.httpServer.ListenAndServeTLS(config.PathCertFile, config.PathKeyFile) //nolint:wrapcheck
	}

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

package server

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"net/http"
	"strings"
	"time"

	categoryrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/mux"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/repository"
	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
)

const (
	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(config *config.Config) error {
	baseCtx := context.Background()

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

	productService, err := usecases.NewProductService(productStorage, *basketService)
	if err != nil {
		return err
	}

	userStorage, err := userrepo.NewUserStorage(pool)
	if err != nil {
		return err
	}

	categoryStorage, err := categoryrepo.NewCategoryStorage(pool)
	if err != nil {
		return err
	}

	handler, err := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin,
		config.Schema, config.PortServer, config.FileServiceDir),
		userStorage, productService, categoryStorage, logger)
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

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

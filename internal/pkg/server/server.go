package server

import (
	"context"
	"net/http"
	"time"

	categoryrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/mux"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/usecases"
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

	logger, err := usecases.NewLogger([]string{"stdout"}, []string{"stderr"})
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer logger.Sync()

	productStorage := productrepo.NewProductStorage(pool, logger)
	userStorage := userrepo.NewUserStorage(pool, logger)
	categoryStorage := categoryrepo.NewCategoryStorage(pool, logger)

	handler := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin, config.Schema, config.PortServer),
		userStorage, productStorage, categoryStorage, logger)

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

package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
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
		return err
	}

	userStorage := userrepo.NewUserStorage(pool)
	productStorage := productrepo.NewProductStorage(pool)

	handler := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin, config.Schema, config.PortServer),
		userStorage, productStorage)

	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.PortServer,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    basicTimeout,
		WriteTimeout:   basicTimeout,
	}

	log.Printf("Start server:%s", config.PortServer)

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

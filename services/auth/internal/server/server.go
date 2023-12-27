package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/interceptors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	reposhare "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/jwt"
	deliverymux "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/usecases"
	"google.golang.org/grpc"
)

const (
	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

// RunFull chErrHTTP - chan from which error can be read if the http server exits with an error
func (s *Server) RunFull(config *config.Config, logger *mylogger.MyLogger, chErrHTTP chan<- error) error { //nolint:funlen
	baseCtx := context.Background()

	handler, err := deliverymux.NewMux(baseCtx, config.AuthServiceName, logger)
	if err != nil {
		return err //nolint:wrapcheck
	}

	go func() {
		s.httpServer = &http.Server{ //nolint:exhaustruct
			Addr:           ":" + config.AuthServicePort,
			Handler:        handler,
			MaxHeaderBytes: http.DefaultMaxHeaderBytes,
			ReadTimeout:    basicTimeout,
			WriteTimeout:   basicTimeout,
		}

		logger.Infof("starting server http at: %s", config.AuthServicePort)

		var err error

		if config.ProductionMode {
			err = s.httpServer.ListenAndServeTLS(config.PathCertFile, config.PathKeyFile)
		} else {
			err = s.httpServer.ListenAndServe()
		}

		chErrHTTP <- err
	}()

	lis, err := net.Listen("tcp", config.AddressAuthServiceGrpc)
	if err != nil {
		return err //nolint:wrapcheck
	}

	metricManager := metrics.NewMetricManagerGrpc(config.AuthServiceName)

	grpcAccessInterceptor := interceptors.NewGrpcAccessInterceptor(metricManager)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcAccessInterceptor.AccessInterceptor, interceptors.ErrConvertInterceptor))

	pool, err := reposhare.NewPgxPool(baseCtx, config.URLDataBase)
	if err != nil {
		return err //nolint:wrapcheck
	}

	storage, err := repository.NewAuthStorage(pool)
	if err != nil {
		return err //nolint:wrapcheck
	}

	service, err := usecases.NewAuthService(storage)
	if err != nil {
		return err //nolint:wrapcheck
	}

	sessionManager, err := delivery.NewSessionManager(pool, service)
	if err != nil {
		return err //nolint:wrapcheck
	}

	auth.RegisterSessionMangerServer(server, sessionManager)

	chCloseRefreshing := make(chan struct{})

	// don`t want use chCloseRefreshing secret now
	err = jwt.StartRefreshingSecret(jwt.TimeTokenLife, chCloseRefreshing)
	if err != nil {
		return err //nolint:wrapcheck
	}

	s.grpcServer = server

	logger.Infof("starting server at: %s", config.AddressAuthServiceGrpc)

	return server.Serve(lis) //nolint:wrapcheck
}

func (s *Server) ShutdownHTTPServer(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

func (s *Server) GracefulStopGrpcServer() {
	s.grpcServer.GracefulStop()
}

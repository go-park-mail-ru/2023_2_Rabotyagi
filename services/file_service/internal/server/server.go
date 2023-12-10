package server

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"net"
	"net/http"
	"strings"
	"time"

	fileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/interceptors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/config"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/delivery/mux"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/usecases"

	"google.golang.org/grpc"
)

const (
	pathCertFile = "/etc/ssl/goods-galaxy.ru.crt"
	pathKeyFile  = "/etc/ssl/goods-galaxy.ru.key"

	urlPrefixPathFS = "img/"

	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

// RunFull chErrHTTP - chan from which error can be read if the http server exits with an error
func (s *Server) RunFull(config *config.Config, chErrHTTP chan<- error) error { //nolint:funlen
	logger, err := my_logger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer logger.Sync() //nolint:errcheck

	baseCtx := context.Background()

	fileStorage, err := repository.NewFileSystemStorage(config.FileServiceDir)
	if err != nil {
		return err //nolint:wrapcheck
	}

	fileServiceHTTP, err := usecases.NewFileServiceHTTP(fileStorage, urlPrefixPathFS)
	if err != nil {
		return err //nolint:wrapcheck
	}

	handler, err := mux.NewMux(baseCtx,
		mux.NewConfigMux(config.AllowOrigin, config.Schema, config.Port, config.FileServiceDir, config.ServiceName),
		fileServiceHTTP, logger)
	if err != nil {
		return err //nolint:wrapcheck
	}

	go func() {
		s.httpServer = &http.Server{ //nolint:exhaustruct
			Addr:           ":" + config.Port,
			Handler:        handler,
			MaxHeaderBytes: http.DefaultMaxHeaderBytes,
			ReadTimeout:    basicTimeout,
			WriteTimeout:   basicTimeout,
		}

		logger.Infof("starting server http at: %s", config.Port)

		var err error

		if config.ProductionMode {
			err = s.httpServer.ListenAndServeTLS(pathCertFile, pathKeyFile)
		} else {
			err = s.httpServer.ListenAndServe()
		}

		chErrHTTP <- err
	}()

	lis, err := net.Listen("tcp", config.AddressFileServiceGrpc)
	if err != nil {
		return err //nolint:wrapcheck
	}

	metricManager := metrics.NewMetricManagerGrpc(config.ServiceName)

	grpcAccessInterceptor := interceptors.NewGrpcAccessInterceptor(metricManager)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcAccessInterceptor.AccessInterceptor, interceptors.ErrConvertInterceptor))

	fileServiceGrpc := usecases.NewFileServiceGrpc(urlPrefixPathFS, fileStorage)
	fileHandlerGrpc := delivery.NewFileHandlerGrpc(fileServiceGrpc)

	fileservice.RegisterFileServiceServer(server, fileHandlerGrpc)

	logger.Infof("starting server grpc at: %s", config.AddressFileServiceGrpc)

	s.grpcServer = server

	return server.Serve(lis) //nolint:wrapcheck
}

func (s *Server) ShutdownHTTPServer(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}

func (s *Server) GracefulStopGrpcServer() {
	s.grpcServer.GracefulStop()
}

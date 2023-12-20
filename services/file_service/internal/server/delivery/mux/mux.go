package mux

import (
	"context"
	"net/http"

	pkgdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/delivery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ConfigMux struct {
	allowOrigin     string
	schema          string
	portServer      string
	fileServiceDir  string
	fileServiceName string
}

func NewConfigMux(allowOrigin string,
	schema string, portServer string, fileServiceDir string, fileServiceName string,
) *ConfigMux {
	return &ConfigMux{
		allowOrigin:     allowOrigin,
		schema:          schema,
		portServer:      portServer,
		fileServiceDir:  fileServiceDir,
		fileServiceName: fileServiceName,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux,
	fileServiceHTTP delivery.IFileServiceHTTP,
	logger *mylogger.MyLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	fileHandler := delivery.NewFileHandlerHTTP(fileServiceHTTP, logger, configMux.fileServiceDir)

	router.Handle("/img/", fileHandler.DocFileServerHandler())
	router.Handle("/img/upload", middleware.Context(ctx,
		middleware.SetupCORS(fileHandler.UploadFileHandler, configMux.allowOrigin, configMux.schema)))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/healthcheck", pkgdelivery.HealthCheckHandler)

	metricsManager := metrics.NewMetricManagerHTTP(configMux.fileServiceName)
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.Context(ctx,
		middleware.AddReqID(
			middleware.AccessLogMiddleware(
				middleware.AddAPIName(router, middleware.APINameV1),
				logger, metricsManager))),
		logger))

	return mux, nil
}

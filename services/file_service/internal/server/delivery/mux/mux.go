package mux

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/delivery"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ConfigMux struct {
	allowOrigin    string
	schema         string
	portServer     string
	fileServiceDir string
}

func NewConfigMux(allowOrigin string, schema string, portServer string, fileServiceDir string) *ConfigMux {
	return &ConfigMux{
		allowOrigin:    allowOrigin,
		schema:         schema,
		portServer:     portServer,
		fileServiceDir: fileServiceDir,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux,
	fileServiceHTTP delivery.IFileServiceHTTP,
	logger *my_logger.MyLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	fileHandler := delivery.NewFileHandlerHTTP(fileServiceHTTP, logger, configMux.fileServiceDir)

	router.Handle("/api/v1/img/", fileHandler.DocFileServerHandler())
	router.Handle("/api/v1/img/upload", middleware.Context(ctx,
		middleware.SetupCORS(fileHandler.UploadFileHandler, configMux.allowOrigin, configMux.schema)))
	router.Handle("/api/v1/metrics", promhttp.Handler())

	metricsManager := metrics.NewMetricManagerHTTP()
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.Context(ctx,
		middleware.AddReqID(middleware.AccessLogMiddleware(router, logger, metricsManager))), logger))

	return mux, nil
}

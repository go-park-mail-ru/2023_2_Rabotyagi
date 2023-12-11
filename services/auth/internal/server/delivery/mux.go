package delivery

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMux(ctx context.Context, fileServiceName string,
	logger *my_logger.MyLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	router.Handle("/api/v1/metrics", promhttp.Handler())

	metricsManager := metrics.NewMetricManagerHTTP(fileServiceName)
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.Context(ctx,
		middleware.AddReqID(middleware.AccessLogMiddleware(router, logger, metricsManager))), logger))

	return mux, nil
}
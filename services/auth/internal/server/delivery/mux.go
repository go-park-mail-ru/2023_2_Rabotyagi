package delivery

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMux(ctx context.Context, fileServiceName string,
	logger *mylogger.MyLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/healthcheck", delivery.HealthCheckHandler)

	metricsManager := metrics.NewMetricManagerHTTP(fileServiceName)
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.Context(ctx,
		middleware.AddReqID(
			middleware.AccessLogMiddleware(
				middleware.AddAPIName(router, middleware.APINameV1),
				logger, metricsManager))),
		logger))

	return mux, nil
}

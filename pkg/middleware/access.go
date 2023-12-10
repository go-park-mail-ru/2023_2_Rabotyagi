package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

type WriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func (w *WriterWithStatus) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func AccessLogMiddleware(next http.Handler,
	logger *my_logger.MyLogger, metricsManager metrics.IMetricManagerHTTP,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writerWithStatus := &WriterWithStatus{ResponseWriter: w, Status: http.StatusOK}

		r = r.WithContext(statuses.FillStatusCtx(r.Context(), statuses.StatusResponseSuccessful))

		start := time.Now()
		next.ServeHTTP(writerWithStatus, r)
		duration := time.Since(start)

		status := statuses.GetStatusCtx(r.Context(), statuses.StatusNotExist) //nolint:contextcheck
		if writerWithStatus.Status != http.StatusOK {
			status = writerWithStatus.Status
		}

		logger := logger.LogReqID(r.Context()) //nolint:contextcheck

		path := r.URL.Path
		method := r.Method
		statusStr := strconv.Itoa(status)

		logger.Infof(
			"path: %s method: %s status: %s duration: %v remoreAddr: %s",
			path, method, statusStr, duration, r.RemoteAddr)

		path = utils.SimplifyPath(path)

		if status == http.StatusNotFound {
			path = "" // this need for reduce count labels in metrics
		}

		metricsManager.IncreaseTotal(path, method, statusStr)
		metricsManager.AddDuration(path, method, statusStr, duration)
	})
}

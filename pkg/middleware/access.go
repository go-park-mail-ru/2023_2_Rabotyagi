package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

type WriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func (w *WriterWithStatus) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func AccessLogMiddleware(next http.Handler, logger *my_logger.MyLogger) http.Handler {
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

		logger.Infof(
			"path: %s method: %s status: %d duration: %v remoreAddr: %s",
			r.URL.Path, r.Method, status, duration, r.RemoteAddr)
	})
}

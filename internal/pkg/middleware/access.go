package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func AccessLogMiddleware(next http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		logger.Infof("%s method: %s RemoteAddr: %s", r.URL.Path, r.Method, r.RemoteAddr)
	})
}

package middleware

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"net/http"
)

func AccessLogMiddleware(next http.Handler, logger *my_logger.MyLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		logger.Infof("%s method: %s RemoteAddr: %s", r.URL.Path, r.Method, r.RemoteAddr)
	})
}

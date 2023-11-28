package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
)

func AddReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := my_logger.AddRequestIDToCtx(r.Context())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

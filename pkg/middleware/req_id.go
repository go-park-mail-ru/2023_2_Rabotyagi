package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"

	"google.golang.org/grpc/metadata"
)

func AddReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := my_logger.AddRequestIDToCtx(r.Context())
		ctx = metadata.NewOutgoingContext(ctx, my_logger.NewMDFromRequestIDCtx(ctx))

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

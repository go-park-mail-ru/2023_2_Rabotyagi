package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"

	"google.golang.org/grpc/metadata"
)

func AddReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := mylogger.AddRequestIDToCtx(r.Context())
		ctx = metadata.NewOutgoingContext(ctx, mylogger.NewMDFromRequestIDCtx(ctx))

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"log"
	"net/http"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %+v\n", err)
				delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

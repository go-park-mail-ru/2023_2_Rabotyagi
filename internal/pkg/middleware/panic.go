package middleware

import (
	"net/http"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			//if err := recover(); err != nil {
			//	log.Println("panic recovered: ", err)
			//	delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))
			//}
		}()
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/statuses"
	"net/http"
)

func Panic(next http.Handler, logger *my_logger.MyLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %+v\n", err)
				delivery.SendResponse(w, logger,
					delivery.NewErrResponse(statuses.StatusInternalServer, delivery.ErrInternalServer))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

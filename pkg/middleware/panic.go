package middleware

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/statuses"
	"net/http"

	"go.uber.org/zap"
)

func Panic(next http.Handler, logger *zap.SugaredLogger) http.Handler {
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

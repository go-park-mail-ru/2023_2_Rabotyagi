package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"

	"go.uber.org/zap"
)

func Panic(next http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %+v\n", err)
				responses.SendResponse(w, logger,
					responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

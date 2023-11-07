package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"

	"go.uber.org/zap"
)

// GetUserIDFromCookie return 0 if error happen and return userID if success
func GetUserIDFromCookie(r *http.Request, logger *zap.SugaredLogger) uint64 {
	cookie, err := r.Cookie(CookieAuthName)
	if err != nil {
		logger.Errorf("in getUserIDFromCookie: %+v\n", err)

		return 0
	}

	rawJwt := cookie.Value

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.Secret)
	if err != nil {
		logger.Errorf("in getUserIDFromCookie: %+v\n", err)

		return 0
	}

	return userPayload.UserID
}

package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/usecases/my_logger"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
)

func GetUserIDFromCookie(r *http.Request) (uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, err
	}

	cookie, err := r.Cookie(CookieAuthName)
	if err != nil {
		logger.Errorln(err)

		return 0, err
	}

	rawJwt := cookie.Value

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.Secret)
	if err != nil {
		logger.Errorln(err)

		return 0, err
	}

	return userPayload.UserID, nil
}

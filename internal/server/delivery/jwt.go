package delivery

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"net/http"
)

func GetUserIDFromCookie(r *http.Request) (uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	cookie, err := r.Cookie(responses.CookieAuthName)
	if err != nil {
		logger.Errorln(err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, responses.ErrCookieNotPresented)
	}

	rawJwt := cookie.Value

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.GetSecret())
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return userPayload.UserID, nil
}

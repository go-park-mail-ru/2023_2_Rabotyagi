package delivery

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
)

func GetUserIDFromCookie(r *http.Request) (uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	cookie, err := r.Cookie(CookieAuthName)
	if err != nil {
		logger.Errorln(err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, ErrCookieNotPresented)
	}

	rawJwt := cookie.Value

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.GetSecret())
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return userPayload.UserID, nil
}

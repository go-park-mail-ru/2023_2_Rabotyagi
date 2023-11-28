package delivery

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"net/http"
)

func GetUserID(ctx context.Context, r *http.Request,
	sessionManager auth.SessionMangerClient) (uint64, error) {
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
	session := auth.Session{AccessToken: rawJwt}

	userID, err := sessionManager.Check(ctx, &session)

	return userID.GetUserId(), nil
}

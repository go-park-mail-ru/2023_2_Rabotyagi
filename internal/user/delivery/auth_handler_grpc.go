package delivery

import (
	"context"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	delivery2 "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	sessionManagerClient auth.SessionMangerClient
	logger               *zap.SugaredLogger
}

func NewAuthHandler(sessionManagerClient auth.SessionMangerClient, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{sessionManagerClient: sessionManagerClient, logger: logger}
}

func (a *AuthHandler) SingUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)

	userWithoutID := new(models.User)
	if err := decoder.Decode(userWithoutID); err != nil {
		delivery2.HandleErr(w, a.logger, err)

		return
	}

	_, err := govalidator.ValidateStruct(userWithoutID)
	if err != nil {
		delivery2.HandleErr(w, a.logger, err)

		return
	}

	userForCreate := auth.User{Email: userWithoutID.Email, Password: userWithoutID.Password}

	sessionWithToken, err := a.sessionManagerClient.Create(ctx, &userForCreate)
	if err != nil {
		delivery2.HandleErr(w, a.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery2.CookieAuthName,
		Value:   sessionWithToken.AccessToken,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery2.SendResponse(w, a.logger, delivery2.NewResponseSuccessful(ResponseSuccessfulSignUp))
	a.logger.Infof("in SignUpHandler: added user")
}

func (a *AuthHandler) SingInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	userForLogin := auth.User{Email: email, Password: password}

	sessionWithToken, err := a.sessionManagerClient.Login(ctx, &userForLogin)
	if err != nil {
		delivery2.HandleErr(w, a.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    delivery2.CookieAuthName,
		Value:   sessionWithToken.AccessToken,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery2.SendResponse(w, a.logger, delivery2.NewResponseSuccessful(ResponseSuccessfulSignUp))
	a.logger.Infof("in SignUpHandler: added user")
}

func (a *AuthHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(delivery2.CookieAuthName)
	if err != nil {
		delivery2.HandleErr(w, a.logger, err)

		return
	}

	sessionUser := &auth.Session{
		AccessToken: cookie.Value,
	}

	ctx := context.Background()
	_, err = a.sessionManagerClient.Delete(ctx, sessionUser)
	if err != nil {
		delivery2.HandleErr(w, a.logger, err)

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery2.SendResponse(w, a.logger, delivery2.NewResponseSuccessful(ResponseSuccessfulLogOut))
	a.logger.Infof("in LogOutHandler: logout user with cookie: %+v", cookie)
}

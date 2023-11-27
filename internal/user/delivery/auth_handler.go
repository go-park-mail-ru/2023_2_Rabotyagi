package delivery

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/statuses"
)

// SignUpHandler godoc
//
//	@Summary    signup
//	@Description  signup in app
//
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags auth
//
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body internal_models.UserWithoutID true  "user data for signup"
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signup [post]
func (u *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	user, err := u.service.AddUser(ctx, r.Body)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.GetSecret())
	if err != nil {
		delivery.SendResponse(w, u.logger,
			delivery.NewErrResponse(statuses.StatusInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:     delivery.CookieAuthName,
		Value:    jwtStr,
		Expires:  expire,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	delivery.SendResponse(w, u.logger, delivery.NewResponseSuccessful(ResponseSuccessfulSignUp))
	u.logger.Infof("in SignUpHandler: added user: %+v", user)
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Tags auth
//	@Produce    json
//	@Param      email  query string true  "user email for signin"
//	@Param      password  query string true  "user password for signin"
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signin [get]
func (u *UserHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	user, err := u.service.GetUser(ctx, email, password)
	if err != nil {
		delivery.HandleErr(w, u.logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(&jwt.UserJwtPayload{
		UserID: user.ID,
		Email:  user.Email,
		Expire: expire.Unix(),
	},
		jwt.GetSecret(),
	)
	if err != nil {
		delivery.SendResponse(w, u.logger,
			delivery.NewErrResponse(statuses.StatusInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:     delivery.CookieAuthName,
		Value:    jwtStr,
		Expires:  expire,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	delivery.SendResponse(w, u.logger, delivery.NewResponseSuccessful(ResponseSuccessfulSignIn))
	u.logger.Infof("in SignInHandler: signin user: %+v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Tags auth
//	@Produce    json
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /logout [post]
func (u *UserHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(delivery.CookieAuthName)
	if err != nil {
		u.logger.Errorln(err)
		delivery.SendResponse(w, u.logger, delivery.NewErrResponse(StatusUnauthorized, ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery.SendResponse(w, u.logger, delivery.NewResponseSuccessful(ResponseSuccessfulLogOut))
	u.logger.Infof("in LogOutHandler: logout user with cookie: %+v", cookie)
}

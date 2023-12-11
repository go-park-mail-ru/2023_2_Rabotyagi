package delivery

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
)

const (
	timeTokenLife = 24 * time.Hour
)

type AuthHandler struct {
	sessionManagerClient auth.SessionMangerClient
	logger               *my_logger.MyLogger
}

func NewAuthHandler(sessionManagerClient auth.SessionMangerClient) (*AuthHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &AuthHandler{sessionManagerClient: sessionManagerClient, logger: logger}, nil
}

// SignUpHandler godoc
//
//	@Summary    signup
//	@Description  signup in app
//
//	@Tags auth
//
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body models.UserWithoutID true  "user data for signup"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /signup [post]
func (a *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := a.logger.LogReqID(ctx)

	decoder := json.NewDecoder(r.Body)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		err = myerrors.NewErrorBadFormatRequest(err.Error())
		responses.HandleErr(w, r, logger, err)

		return
	}

	_, err := govalidator.ValidateStruct(userWithoutID)
	if err != nil {
		err = myerrors.NewErrorBadContentRequest(err.Error())
		responses.HandleErr(w, r, logger, err)

		return
	}

	userForCreate := auth.User{Email: userWithoutID.Email, Password: userWithoutID.Password}

	sessionWithToken, err := a.sessionManagerClient.Create(ctx, &userForCreate)
	if err != nil {
		logger.Errorln(err)
		responses.HandleErr(w, r, logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    responses.CookieAuthName,
		Value:   sessionWithToken.GetAccessToken(),
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	responses.SendResponse(w, logger, responses.NewResponseSuccessful(ResponseSuccessfulSignUp))
	logger.Infof("in SignUpHandler: added user")
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Tags auth
//	@Produce    json
//	@Param      email  query string true  "user email for signin"
//	@Param      password  query string true  "user password for signin"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /signin [get]
func (a *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := a.logger.LogReqID(ctx)

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	userForLogin := auth.User{Email: email, Password: password}

	sessionWithToken, err := a.sessionManagerClient.Login(ctx, &userForLogin)
	if err != nil {
		logger.Errorln(err)
		responses.HandleErr(w, r, logger, err)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    responses.CookieAuthName,
		Value:   sessionWithToken.GetAccessToken(),
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	responses.SendResponse(w, logger, responses.NewResponseSuccessful(ResponseSuccessfulSignIn))
	logger.Infof("in SignInHandler: added user")
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Tags auth
//	@Produce    json
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error".
//	@Router      /logout [post]
func (a *AuthHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := a.logger.LogReqID(ctx)

	cookie, err := r.Cookie(responses.CookieAuthName)
	if err != nil {
		logger.Errorln(err)

		if errors.Is(err, http.ErrNoCookie) {
			err = responses.ErrCookieNotPresented
		}

		responses.HandleErr(w, r, logger, err)

		return
	}

	sessionUser := &auth.Session{
		AccessToken: cookie.Value,
	}

	expiredSession, err := a.sessionManagerClient.Delete(ctx, sessionUser)
	if err != nil {
		logger.Errorln(err)
		responses.HandleErr(w, r, logger, err)

		return
	}

	cookie.Value = expiredSession.GetAccessToken()

	http.SetCookie(w, cookie)
	responses.SendResponse(w, logger, responses.NewResponseSuccessful(ResponseSuccessfulLogOut))
	logger.Infof("in LogOutHandler: logout user with cookie: %+v", cookie)
}

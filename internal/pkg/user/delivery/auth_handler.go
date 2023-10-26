package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
)

const (
	CookieAuthName = "access_token"
	timeTokenLife  = 24 * time.Hour
)

type AuthHandler struct {
	Storage    usecases.IUserStorage
	AddrOrigin string
}

// SignUpHandler godoc
//
//	@Summary    signup
//	@Description  signup in app
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body models.PreUser true  "user data for signup"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signup [post]
func (a *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, a.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		log.Printf("error in SignUpHandler: %v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, delivery.ErrBadRequest))

		return
	}

	user, err := a.Storage.AddUser(ctx, userWithoutID)
	if err != nil {
		log.Printf("error in SignUpHandler: %v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer)) //TODO maybe not at all 500

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.Secret)
	if err != nil {
		log.Printf("%v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignUp))
	log.Printf("in SignUpHandler: added user: %v", user)
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body models.PreUser true  "user data for signin"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /signin [post]
func (a *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, a.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	userWithoutID := new(models.UserWithoutID)

	if err := decoder.Decode(userWithoutID); err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, delivery.ErrBadRequest))

		return
	}

	exist, err := a.Storage.IsUserExist(ctx, userWithoutID.Email, userWithoutID.Phone)
	if !exist {
		log.Printf("user is not exists %v\n", userWithoutID)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrUserNotExits))

		return
	}

	if err != nil {
		log.Printf("%v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	user, err := a.Storage.GetUserByEmail(ctx, userWithoutID.Email)
	if err != nil || userWithoutID.Pass != user.Pass {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrWrongCredentials))

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(&jwt.UserJwtPayload{
		UserID: user.ID,
		Email:  user.Email,
		Expire: expire.Unix(),
	},
		jwt.Secret,
	)
	if err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulSignIn))
	log.Printf("signin user: %v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /logout [post]
func (a *AuthHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, a.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(CookieAuthName)
	if err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(StatusUnauthorized, ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulLogOut))
	log.Printf("logout user with cookie: %v", cookie)
}

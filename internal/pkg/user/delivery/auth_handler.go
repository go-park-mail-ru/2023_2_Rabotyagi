package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
)

const (
	CookieAuthName = "access_token"
	timeTokenLife  = 24 * time.Hour
)

type AuthHandler struct {
	Storage    *repository.AuthStorageMap
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
func (h *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, h.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(models.PreUser)
	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, delivery.ErrBadRequest))

		return
	}

	if h.Storage.IsUserExist(preUser.Email) {
		log.Printf("already exist user %v\n", preUser)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrUserAlreadyExist))

		return
	}

	err := h.Storage.CreateUser(preUser)
	if err != nil {
		log.Printf("%v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	user, err := h.Storage.GetUser(preUser.Email)
	if err != nil {
		log.Printf("%v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

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
	log.Printf("added user: %v", user)
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
func (h *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, h.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(models.PreUser)

	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, delivery.ErrBadRequest))

		return
	}

	if !h.Storage.IsUserExist(preUser.Email) {
		log.Printf("user is not exists %v\n", preUser)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrUserNotExits))

		return
	}

	user, err := h.Storage.GetUser(preUser.Email)
	if err != nil || preUser.Password != user.Password {
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
func (h *AuthHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, h.AddrOrigin)

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

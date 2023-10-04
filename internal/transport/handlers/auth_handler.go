package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
	resp "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/responses"
)

const timeTokenLife = 24 * time.Hour

// signUpHandler godoc
//
//	@Summary    signup
//	@Description  signup in app
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body storage.PreUser true  "user data for signup"
//	@Success    200  {object} Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} ErrorResponse "Error"
//	@Router      /signup [post]
func (h *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp.SetupCORS(&w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(storage.PreUser)
	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrBadRequest))

		return
	}

	if h.Storage.IsUserExist(preUser.Email) {
		log.Printf("already exist user %v\n", preUser)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrUserAlreadyExist))

		return
	}

	err := h.Storage.CreateUser(preUser)
	if err != nil {
		log.Printf("%v", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrInternalServer, resp.ErrInternalServer))

		return
	}

	user, err := h.Storage.GetUser(preUser.Email)
	if err != nil {
		log.Printf("%v", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrInternalServer, resp.ErrInternalServer))

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.Secret)
	if err != nil {
		log.Printf("%v", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrInternalServer, resp.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    resp.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
	}

	http.SetCookie(w, cookie)
	resp.SendOkResponse(w, resp.NewResponse(resp.StatusResponseSuccessful, resp.ResponseSuccessfulSignUp))
	log.Printf("added user: %v", user)
}

// SignInHandler godoc
//
//	@Summary    signin
//	@Description  signin in app
//	@Accept      json
//	@Produce    json
//	@Param      preUser  body storage.PreUser true  "user data for signin"
//	@Success    200  {object} Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} ErrorResponse "Error"
//	@Router      /signin [post]
func (h *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp.SetupCORS(&w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(storage.PreUser)

	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrBadRequest))

		return
	}

	if !h.Storage.IsUserExist(preUser.Email) {
		log.Printf("user is not exists %v\n", preUser)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrUserNotExits))

		return
	}

	user, err := h.Storage.GetUser(preUser.Email)
	if err != nil || preUser.Password != user.Password {
		log.Printf("%v\n", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrWrongCredentials))

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
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrInternalServer, resp.ErrInternalServer))

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    resp.CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	resp.SendOkResponse(w, resp.NewResponse(resp.StatusResponseSuccessful, resp.ResponseSuccessfulSignIn))
	log.Printf("signin user: %v", user)
}

// LogOutHandler godoc
//
//	@Summary    logout
//	@Description  logout in app
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} ErrorResponse "Error"
//	@Router      /logout [post]
func (h *AuthHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp.SetupCORS(&w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(resp.CookieAuthName)
	if err != nil {
		log.Printf("%v\n", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusUnauthorized, resp.ErrUnauthorized))

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	resp.SendOkResponse(w, resp.NewResponse(resp.StatusResponseSuccessful, resp.ResponseSuccessfulLogOut))
	log.Printf("logout user with cookie: %v", cookie)
}

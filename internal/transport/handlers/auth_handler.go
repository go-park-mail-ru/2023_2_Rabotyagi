package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

const timeTokenLife = 24 * time.Hour

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

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
//	@Failure    200  {object} ErrorResponse
//	@Router      /signup [post]
func (h *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setupCORS(&w, r)

	if (*r).Method == "OPTIONS" {
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
		sendResponse(w, ErrBadRequest)

		return
	}

	if h.Storage.IsUserExist(preUser.Email) {
		log.Printf("already exist user %v\n", preUser)
		sendResponse(w, ErrUserAlreadyExist)

		return
	}

	err := h.Storage.CreateUser(preUser)
	if err != nil {
		log.Printf("%v", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	user, err := h.Storage.GetUser(preUser.Email)
	if err != nil {
		log.Printf("%v", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	expire := time.Now().Add(timeTokenLife)

	jwtStr, err := jwt.GenerateJwtToken(
		&jwt.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: expire.Unix()}, jwt.Secret)
	if err != nil {
		log.Printf("%v", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
	}

	http.SetCookie(w, cookie)
	sendResponse(w, ResponseSuccessfulSignUp)
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
//	@Failure    200  {object} ErrorResponse
//	@Router      /signin [post]
func (h *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setupCORS(&w, r)

	if (*r).Method == "OPTIONS" {
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
		sendResponse(w, ErrBadRequest)

		return
	}

	if !h.Storage.IsUserExist(preUser.Email) {
		log.Printf("user is not exists %v\n", preUser)
		sendResponse(w, ErrWrongCredentials)

		return
	}

	user, err := h.Storage.GetUser(preUser.Email)
	if err != nil || preUser.Password != user.Password {
		log.Printf("%v\n", err)
		sendResponse(w, ErrWrongCredentials)

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
		sendResponse(w, ErrInternalServer)

		return
	}

	cookie := &http.Cookie{ //nolint:exhaustruct
		Name:    CookieAuthName,
		Value:   jwtStr,
		Expires: expire,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	sendResponse(w, ResponseSuccessfulSignIn)
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
//	@Failure    200  {object} ErrorResponse
//	@Router      /logout [post]
func (h *AuthHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	setupCORS(&w, r)

	if (*r).Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	cookie, err := r.Cookie(CookieAuthName)
	if err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrUnauthorized)

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	sendResponse(w, ResponseSuccessfulLogOut)
	log.Printf("logout user with cookie: %v", cookie)
}

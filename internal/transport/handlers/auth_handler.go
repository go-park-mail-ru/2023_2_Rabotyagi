package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

const timeTokenLife = 24 * time.Hour

func (h *AuthHandler) signUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(storage.PreUser)
	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}

	if h.storage.IsUserExist(preUser.Email) {
		log.Printf("already exist user %v\n", preUser)
		sendResponse(w, ErrUserAlreadyExist)

		return
	}

	err := h.storage.CreateUser(preUser)
	if err != nil {
		log.Printf("%v", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	user, err := h.storage.GetUser(preUser.Email)
	if err != nil {
		log.Printf("%v", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	jwtStr, err := auth.GenerateJwtToken(
		&auth.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: time.Now().Add(timeTokenLife).Unix()})
	if err != nil {
		log.Printf("%v", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	cookie := &http.Cookie{ //nolint:exhaustruct,exhaustivestruct
		Name:  CookieAuthName,
		Value: jwtStr,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulSignUp)

	log.Printf("added user: %v", user)
}

func (h *AuthHandler) signInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(storage.PreUser)

	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}

	if !h.storage.IsUserExist(preUser.Email) {
		log.Printf("user is not exists %v\n", preUser)
		sendResponse(w, ErrWrongCredentials)

		return
	}

	user, err := h.storage.GetUser(preUser.Email)
	if err != nil || preUser.Password != user.Password {
		log.Printf("%v\n", err)
		sendResponse(w, ErrWrongCredentials)

		return
	}

	jwtStr, err := auth.GenerateJwtToken(
		&auth.UserJwtPayload{UserID: user.ID, Email: user.Email, Expire: time.Now().Add(timeTokenLife).Unix()})
	if err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrInternalServer)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	cookie := &http.Cookie{ //nolint:exhaustruct,exhaustivestruct
		Name:  CookieAuthName,
		Value: jwtStr,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulSignIn)

	log.Printf("sign user: %v", user)
}

func (h *AuthHandler) logOut(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	cookie, err := r.Cookie(CookieAuthName)
	if err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrUnauthorized)

		return
	}

	cookie.Expires = time.Now()

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulLogOut)

	log.Printf("logout user with cookie: %v", cookie)
}

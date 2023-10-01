package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	auth "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/authorization"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

// TODO from config
var secret = []byte("super-secret")

// signUpHandler godoc
//
//  @Summary    register user
//  @Description  register new user
//  @Accept      json
//  @Produce    json
//  @Param      user  body    preUser  true  "User"
//  @Success    200  {object}  myError
//  @Failure    500  {object}  myError
//  @Router      /signup/ [post]
func (h *AuthHandler) signUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)

	preUser := new(storage.PreUser)
	if err := decoder.Decode(preUser); err != nil {
		log.Printf("%v\n", err)
		sendErr(w, ErrBadRequest)

		return
	}

	if h.storage.IsUserExist(preUser.Email) {
		log.Printf("already exist user %v\n", preUser)
		sendErr(w, ErrUserAlreadyExist)

		return
	}

	err := h.storage.CreateUser(preUser)
	if err != nil {
		log.Printf("%v", err)
		sendErr(w, ErrInternalServer)

		return
	}

	user, err := h.storage.GetUser(preUser.Email)
	if err != nil {
		log.Printf("%v", err)
		sendErr(w, ErrInternalServer)

		return
	}

	jwtStr, err := auth.GenerateJwtToken(user, secret)
	if err != nil {
		log.Printf("%v", err)
		sendErr(w, ErrInternalServer)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	cookie := &http.Cookie{ //nolint:exhaustruct,exhaustivestruct
		Name:  "session_id",
		Value: jwtStr,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulSignUp)

	log.Printf("added user: %v", user)
}

var (
	errBadCredentials = errors.New("email or password is incorrect")
)

// signInHandler godoc
//
//  @Summary    login user
//  @Description  login user
//  @Accept      json
//  @Produce    json
//  @Param      user  body    preUser  true  "User"
//  @Success    200  {object}  myError
//  @Failure    500  {object}  myError
//  @Router      /signin/ [post]
func (h *AuthHandler) signInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
  
	if r.Method != http.MethodPost {
	  http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
  
	decoder := json.NewDecoder(r.Body)
  
	preUser := new(storage.PreUser)
  
	if err := decoder.Decode(preUser); err != nil {
	  log.Printf("%v\n", err)
	  sendErr(w, ErrBadRequest)
  
	  return
	}
  
	if !h.storage.IsUserExist(preUser.Email) {
	  log.Printf("user is not exists %v\n", preUser)
	  sendErr(w, ErrWrongCredentials)
  
	  return
	}
  
	user, err := h.storage.GetUser(preUser.Email)
	if err != nil || preUser.Password != user.Password {
	  log.Printf("%v\n", err)
	  sendErr(w, ErrWrongCredentials)
  
	  return
	}
  
	jwtStr, err := auth.GenerateJwtToken(user, secret)
	if err != nil {
	  log.Printf("%v\n", err)
	  sendErr(w, ErrInternalServer)
  
	  return
	}
  
	w.Header().Set("Content-Type", "application/json")
  
	cookie := &http.Cookie{ //nolint:exhaustruct,exhaustivestruct
	  Name:  "session_id",
	  Value: jwtStr,
	}
  
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
  
	sendResponse(w, ResponseSuccessfulSignIn)
  
	log.Printf("sign user: %v", user)
  }

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

type AuthHandler struct {
	storage *storage.AuthStorageMap
}

func (h *AuthHandler) InitRoutes() http.Handler {
	router := http.NewServeMux()

	storageMap := storage.NewAuthStorageMap()
	authHandler := &AuthHandler{
		storage: storageMap,
	}

	router.HandleFunc("/api/v1/signup/", authHandler.signUpHandler)

	return router
}

func sendErr(w http.ResponseWriter, errResponse ErrorResponse) {
	response, err := json.Marshal(errResponse)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}

	_, err = w.Write(response)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}
}

func sendResponse(w http.ResponseWriter, response Response) {
	responseSend, err := json.Marshal(response)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}

	_, err = w.Write(responseSend)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}
}

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

func (h *AuthHandler) signInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// декодим json из реквеста в storage.PreUser
	decoder := json.NewDecoder(r.Body)

	user := new(storage.PreUser)
	err := decoder.Decode(user)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	// нет юзера с таким именем
	if !h.storage.IsUserExist(user.Email) {
		log.Printf("user is not exists")
		w.Write([]byte("{}"))
		return
	}

	// неправильный пароль
	userWithId, err := h.storage.GetUser(user.Email)
	if err != nil || user.Password != userWithId.Password {
		log.Printf("error while getting user: %s", err)
		w.Write([]byte("{}"))
		return
	}

	// генерируем jwt токен из юзера
	jwtStr, err := auth.GenerateJwtToken(userWithId, secret)
	if err != nil {
		log.Printf("%s", err)
		w.Write([]byte("{}"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// выставляем куку
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: jwtStr,
	}

	http.SetCookie(w, cookie)

	// Отправляем токен как строку в теле ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtStr))
}

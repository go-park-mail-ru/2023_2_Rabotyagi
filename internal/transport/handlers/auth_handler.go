package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	auth "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/authorization"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

type AuthHandler struct {
	storage *storage.AuthStorageMap
}

func (h *AuthHandler) InitRoutes() http.Handler {
	router := http.NewServeMux()

	storageMap := storage.NewAuthStorageMap()
	authHandler := &AuthHandler{
		storage: storageMap,
	}

	router.HandleFunc("/api/v1/signin/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Вы на логине")
		log.Println(r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			authHandler.signInHandler(w, r)
		} else {
			http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/api/v1/signup/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			authHandler.signUpHandler(w, r)
		} else {
			http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
		}
	})

	return router
}

var errUserExists = errors.New("the user already exists")

var secret = []byte("super-secret")

func (h *AuthHandler) signUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	newUser := new(storage.PreUser)
	err := decoder.Decode(newUser)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	fmt.Println(newUser)

	// уже есть юзер с таким именем
	if h.storage.IsUserExist(newUser.Email) {
		log.Printf("%s", errUserExists)
		w.Write([]byte("{}"))
		return
	}
	
	// создаем юзера
	h.storage.CreateUser(newUser)

	userWithId, err := h.storage.GetUser(newUser.Email)
	if err != nil {
		log.Printf("error while getting user: %s", err)
		w.Write([]byte("{}"))
		return
	}

	// генерируем jwt токен
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

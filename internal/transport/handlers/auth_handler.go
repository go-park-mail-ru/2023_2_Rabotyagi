package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	rabotyagi "github.com/go-park-mail-ru/2023_2_Rabotyagi"
	auth "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/authorization"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

type AuthHandler struct {
	storage *storage.AuthStorageMap
}

// type RegRequest struct {
// 	Name     string `json:"name"`
// 	Password string `json:"password"`
// }

func (h *AuthHandler) InitRoutes() http.Handler {
	router := http.NewServeMux()
	authHandler := AuthHandler{}

	router.HandleFunc("/signin/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			authHandler.signInHandler(w, r)
		default:
			http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/signup/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			authHandler.signUpHandler(w, r)
		default:
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

	// var id int = 0
	// if len(h.storage.users) > 0 {
	// 	id = h.storage.users[len(h.storage.users)-1].Id + 1
	// }

	if h.storage.IsUserExist(newUser.Name) {
		log.Printf("%s", errUserExists)
		w.Write([]byte("{}"))
		return
	}

	h.storage.CreateUser(newUser)

	userWithId, err := h.storage.GetUser(newUser.Name)
	if err != nil {
		log.Printf("error while getting user: %s", err)
		w.Write([]byte("{}"))
		return
	}

	jwtStr, err := auth.GenerateJwtToken(userWithId, secret)
	if err != nil {
		log.Printf("%s", err)
		w.Write([]byte("{}"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

var jwtSecretKey = []byte("very-secret-key")

func (h *AuthHandler) signInHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	user := new(rabotyagi.User)
	err := decoder.Decode(user)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	if !h.storage.IsUserExist(user.Name) {
		log.Printf("user is not exists")
		w.Write([]byte("{}"))
		return
	}

	userWithId, err := h.storage.GetUser(user.Name)
	if err != nil || user.Password != userWithId.Password {
		log.Printf("error while getting user: %s", err)
		w.Write([]byte("{}"))
		return
	}

	jwtStr, err := auth.GenerateJwtToken(userWithId, secret)
	if err != nil {
		log.Printf("%s", err)
		w.Write([]byte("{}"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: jwtStr,
	}

	http.SetCookie(w, cookie)

	// Отправляем токен как строку в теле ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtStr))
}

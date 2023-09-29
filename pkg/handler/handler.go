package handler

import (
    "net/http"
)

type Handler struct {

}

func (h *Handler) InitRoutes() http.Handler {
    router := http.NewServeMux()

    router.HandleFunc("/", homeHandler)
    router.HandleFunc("/sign-up", signUpHandler)
    router.HandleFunc("/sign-in", signInHandler)

	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("дефолтная страница"))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        w.Write([]byte("регистрация нового пользователя"))
    case http.MethodPost:
        w.Write([]byte("обработка регистрации пользователя"))
    default:
        http.Error(w, "пупупууу", http.StatusMethodNotAllowed)
    }
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        w.Write([]byte("аутентификация пользователя"))
    case http.MethodPost:
        w.Write([]byte("обработка аутентификации пользователя"))
    default:
        http.Error(w, "пупупууу", http.StatusMethodNotAllowed)
    }
}
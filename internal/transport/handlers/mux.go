//  @title      YOULA roject API
//  @version    1.0
//  @description  This is a sample server MY server.

// @host    127.0.0.1:8080
// @BasePath  /api/v1

package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

type AuthHandler struct {
	Storage *storage.AuthStorageMap
}

type PostHandler struct {
	Storage *storage.PostStorageMap
}

type Handler struct {
	AuthHandler *AuthHandler
	PostHandler *PostHandler
}

func (h *Handler) InitRoutes() http.Handler {
	router := http.NewServeMux()

	authStorageMap := storage.NewAuthStorageMap()
	authHandler := &AuthHandler {
		Storage: authStorageMap,
	}

	postStorageMap := storage.NewPostStorageMap()
	postHandler := &PostHandler {
		Storage: postStorageMap,
	}

	router.HandleFunc("/api/v1/signup", authHandler.SignUpHandler)
	router.HandleFunc("/api/v1/signin", authHandler.SignInHandler)
	router.HandleFunc("/api/v1/logout", authHandler.LogOutHandler)

	router.HandleFunc("/api/v1/post/add", postHandler.AddPostHandler)
	router.HandleFunc("/api/v1/post/get/", postHandler.GetPostHandler)
	router.HandleFunc("/api/v1/post/get_list", postHandler.GetPostsListHandler)

	return router
}

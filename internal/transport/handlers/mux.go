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
	storage *storage.AuthStorageMap
}

type PostHandler struct {
	storage *storage.PostStorageMap
}

type Handler struct {
    AuthHandler *AuthHandler
    PostHandler *PostHandler
}

func (h *Handler) InitRoutes() http.Handler {
	router := http.NewServeMux()

	authStorageMap := storage.NewAuthStorageMap()
	authHandler := &AuthHandler{
		storage: authStorageMap,
	}

	postStorageMap := storage.NewPostStorageMap()
	postHandler := &PostHandler{
		storage: postStorageMap,
	}	

	router.HandleFunc("/api/v1/signup/", authHandler.signUpHandler)
	router.HandleFunc("/api/v1/signin/", authHandler.signInHandler)

	router.HandleFunc("/api/v1/post/add", postHandler.addPostHandler)
	router.HandleFunc("/api/v1/post/get", postHandler.getPostHandler)
	router.HandleFunc("/api/v1/post/get_list", postHandler.getPostsListHandler)

	return router
}
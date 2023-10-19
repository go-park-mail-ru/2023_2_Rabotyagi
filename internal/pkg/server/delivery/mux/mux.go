package mux

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/middleware"
	"net/http"

	postdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/post/delivery"
	postrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/post/repository"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
)

type Handler struct {
	AuthHandler *userdelivery.AuthHandler
	PostHandler *postdelivery.PostHandler
}

func NewMux(addrOrigin string) http.Handler {
	router := http.NewServeMux()

	authStorageMap := userrepo.NewAuthStorageMap()

	authHandler := &userdelivery.AuthHandler{
		Storage:    authStorageMap,
		AddrOrigin: addrOrigin,
	}

	postStorageMap := postrepo.NewPostStorageMap()
	postHandler := &postdelivery.PostHandler{
		Storage:    postrepo.GeneratePosts(postStorageMap),
		AddrOrigin: addrOrigin,
	}

	router.HandleFunc("/api/v1/signup", authHandler.SignUpHandler)
	router.HandleFunc("/api/v1/signin", authHandler.SignInHandler)
	router.HandleFunc("/api/v1/logout", authHandler.LogOutHandler)

	router.HandleFunc("/api/v1/post/add", postHandler.AddPostHandler)
	router.HandleFunc("/api/v1/post/get/", postHandler.GetPostHandler)
	router.HandleFunc("/api/v1/post/get_list", postHandler.GetPostsListHandler)

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router))

	return mux
}

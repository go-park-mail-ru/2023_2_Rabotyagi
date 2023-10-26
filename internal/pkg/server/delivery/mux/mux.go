package mux

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/middleware"
	postdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/post/delivery"
	postrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/post/repository"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
)

type Handler struct {
	AuthHandler *userdelivery.AuthHandler
	PostHandler *postdelivery.PostHandler
}

func NewMux(ctx context.Context, addrOrigin string, userStorage userrepo.IUserStorage) http.Handler {
	router := http.NewServeMux()

	authHandler := &userdelivery.AuthHandler{
		Storage:    userStorage,
		AddrOrigin: addrOrigin,
	}

	postStorageMap := postrepo.NewPostStorageMap()
	postHandler := &postdelivery.PostHandler{
		Storage:    postrepo.GeneratePosts(postStorageMap),
		AddrOrigin: addrOrigin,
	}

	router.Handle("/api/v1/signup", middleware.Context(ctx, authHandler.SignUpHandler))
	router.Handle("/api/v1/signin", middleware.Context(ctx, authHandler.SignInHandler))
	router.Handle("/api/v1/logout", middleware.Context(ctx, authHandler.LogOutHandler))

	router.Handle("/api/v1/post/add", middleware.Context(ctx, postHandler.AddPostHandler))
	router.Handle("/api/v1/post/get/", middleware.Context(ctx, postHandler.GetPostHandler))
	router.Handle("/api/v1/post/get_list", middleware.Context(ctx, postHandler.GetPostsListHandler))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router))

	return mux
}

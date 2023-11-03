package mux

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/middleware"
	postdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/delivery"
	postrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
)

type Handler struct {
	UserHandler *userdelivery.UserHandler
	PostHandler *postdelivery.PostHandler
}

func NewMux(ctx context.Context, addrOrigin string, userStorage userusecases.IUserStorage) http.Handler {
	router := http.NewServeMux()

	userHandler := &userdelivery.UserHandler{
		Storage:    userStorage,
		AddrOrigin: addrOrigin,
	}

	postStorageMap := postrepo.NewProductStorageMap()
	postHandler := &postdelivery.PostHandler{
		Storage:    postrepo.GenerateProducts(postStorageMap),
		AddrOrigin: addrOrigin,
	}

	imgHandler := http.StripPrefix(
		"/api/v1/img/",
		http.FileServer(http.Dir("./db/img")),
	)

	router.Handle("/api/v1/img/", imgHandler)

	router.Handle("/api/v1/signup", middleware.Context(ctx, userHandler.SignUpHandler))
	router.Handle("/api/v1/signin", middleware.Context(ctx, userHandler.SignInHandler))
	router.Handle("/api/v1/logout", middleware.Context(ctx, userHandler.LogOutHandler))

	router.Handle("/api/v1/profile/get/", middleware.Context(ctx, userHandler.GetUserHandler))
	router.Handle("/api/v1/profile/rebuild", middleware.Context(ctx, userHandler.FullyUpdateUserHandler))
	router.Handle("/api/v1/profile/update", middleware.Context(ctx, userHandler.PartiallyUpdateUserHandler))

	router.Handle("/api/v1/post/add", middleware.Context(ctx, postHandler.AddPostHandler))
	router.Handle("/api/v1/post/get/", middleware.Context(ctx, postHandler.GetPostHandler))
	router.Handle("/api/v1/post/get_list", middleware.Context(ctx, postHandler.GetPostsListHandler))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router))

	return mux
}

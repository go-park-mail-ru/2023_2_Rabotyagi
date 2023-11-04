package mux

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/middleware"
	productdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/delivery"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
)

type Handler struct {
	AuthHandler *userdelivery.AuthHandler
	PostHandler *productdelivery.PostHandler
}

func NewMux(ctx context.Context, addrOrigin string, userStorage userusecases.IUserStorage, productStorage productusecases.IProductStorage) http.Handler {
	router := http.NewServeMux()

	authHandler := &userdelivery.AuthHandler{
		Storage:    userStorage,
		AddrOrigin: addrOrigin,
	}

	postHandler := &productdelivery.PostHandler{
		Storage:    productStorage,
		AddrOrigin: addrOrigin,
	}

	imgHandler := http.StripPrefix(
		"/api/v1/img/",
		http.FileServer(http.Dir("./db/img")),
	)

	router.Handle("/api/v1/img/", imgHandler)

	router.Handle("/api/v1/signup", middleware.Context(ctx, authHandler.SignUpHandler))
	router.Handle("/api/v1/signin", middleware.Context(ctx, authHandler.SignInHandler))
	router.Handle("/api/v1/logout", middleware.Context(ctx, authHandler.LogOutHandler))

	router.Handle("/api/v1/post/add", middleware.Context(ctx, postHandler.AddPostHandler))
	router.Handle("/api/v1/post/get/", middleware.Context(ctx, postHandler.GetPostHandler))
	//router.Handle("/api/v1/post/get_list", middleware.Context(ctx, postHandler.GetPostsListHandler))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router))

	return mux
}

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

type ConfigMux struct {
	addrOrigin string
	schema     string
	portServer string
}

func NewConfigMux(addrOrigin string, schema string, portServer string) *ConfigMux {
	return &ConfigMux{
		addrOrigin: addrOrigin,
		schema:     schema,
		portServer: portServer,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux, userStorage userusecases.IUserStorage,
	productStorage productusecases.IProductStorage,
) http.Handler {
	router := http.NewServeMux()

	userHandler := userdelivery.NewUserHandler(userStorage, configMux.addrOrigin, configMux.schema)

	productHandler := productdelivery.NewProductHandler(productStorage,
		configMux.addrOrigin, configMux.schema, configMux.portServer,
	)

	imgHandler := http.StripPrefix(
		"/api/v1/img/",
		http.FileServer(http.Dir("./db/img")),
	)

	router.Handle("/api/v1/img/", imgHandler)

	router.Handle("/api/v1/signup", middleware.Context(ctx, userHandler.SignUpHandler))
	router.Handle("/api/v1/signin", middleware.Context(ctx, userHandler.SignInHandler))
	router.Handle("/api/v1/logout", middleware.Context(ctx, userHandler.LogOutHandler))

	router.Handle("/api/v1/profile/get/", middleware.Context(ctx, userHandler.GetUserHandler))
	router.Handle("/api/v1/profile/update", middleware.Context(ctx, userHandler.PartiallyUpdateUserHandler))

	router.Handle("/api/v1/product/add", middleware.Context(ctx, productHandler.AddProductHandler))
	router.Handle("/api/v1/product/get/", middleware.Context(ctx, productHandler.GetProductHandler))
	router.Handle("/api/v1/product/get_list", middleware.Context(ctx, productHandler.GetProductListHandler))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router))

	return mux
}

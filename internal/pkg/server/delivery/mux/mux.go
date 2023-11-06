package mux

import (
	"context"
	"net/http"

	filedelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/delivery"
	filerepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/repository"
	fileusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/middleware"
	productdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/delivery"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
)

const baseDir = "./db/img/"

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

	fileStorage := filerepo.NewFileSystemStorage(baseDir)
	fileService := fileusecases.NewFileService(fileStorage)
	fileHandler := filedelivery.NewFileHandler(baseDir, fileService)

	router.Handle("/api/v1/img/", fileHandler)
	router.Handle("/api/v1/img/upload", middleware.Context(ctx, fileHandler.UploadFileHandler))

	router.Handle("/api/v1/signup", middleware.Context(ctx, userHandler.SignUpHandler))
	router.Handle("/api/v1/signin", middleware.Context(ctx, userHandler.SignInHandler))
	router.Handle("/api/v1/logout", middleware.Context(ctx, userHandler.LogOutHandler))

	router.Handle("/api/v1/profile/get/", middleware.Context(ctx, userHandler.GetUserHandler))
	router.Handle("/api/v1/profile/update", middleware.Context(ctx, userHandler.PartiallyUpdateUserHandler))

	router.Handle("/api/v1/product/add", middleware.Context(ctx, productHandler.AddProductHandler))
	router.Handle("/api/v1/product/get/", middleware.Context(ctx, productHandler.GetProductHandler))
	router.Handle("/api/v1/product/get_list", middleware.Context(ctx, productHandler.GetProductListHandler))
	router.Handle("/api/v1/product/get_list_of_saler",
		middleware.Context(ctx, productHandler.GetListProductOfSalerHandler))
	router.Handle("/api/v1/product/update/", middleware.Context(ctx, productHandler.UpdateProductHandler))
	router.Handle("/api/v1/product/close/", middleware.Context(ctx, productHandler.CloseProductHandler))
	router.Handle("/api/v1/product/delete/", middleware.Context(ctx, productHandler.DeleteProductHandler))

	router.Handle("/api/v1/order/add", middleware.Context(ctx, productHandler.AddOrderHandler))
	router.Handle("/api/v1/order/get_basket", middleware.Context(ctx, productHandler.GetBasketHandler))
	router.Handle("/api/v1/order/update_count", middleware.Context(ctx, productHandler.UpdateOrderCountHandler))
	router.Handle("/api/v1/order/update_status", middleware.Context(ctx, productHandler.UpdateOrderStatusHandler))
	router.Handle("/api/v1/order/buy_full_basket", middleware.Context(ctx, productHandler.BuyFullBasketHandler))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router))

	return mux
}

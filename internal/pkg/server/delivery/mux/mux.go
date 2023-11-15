package mux

import (
	"context"
	"net/http"

	categorydelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/delivery"
	categoryusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/usecases"
	filedelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/delivery"
	filerepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/repository"
	fileusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/middleware"
	productdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/delivery"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"

	"go.uber.org/zap"
)

const urlPrefixPathFS = "img/"

type ConfigMux struct {
	addrOrigin     string
	schema         string
	portServer     string
	fileServiceDir string
}

func NewConfigMux(addrOrigin string, schema string, portServer string, fileServiceDir string) *ConfigMux {
	return &ConfigMux{
		addrOrigin:     addrOrigin,
		schema:         schema,
		portServer:     portServer,
		fileServiceDir: fileServiceDir,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux, userStorage userusecases.IUserStorage,
	productStorage productusecases.IProductStorage, categoryStorage categoryusecases.ICategoryStorage,
	logger *zap.SugaredLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	userHandler, err := userdelivery.NewUserHandler(userStorage, configMux.addrOrigin, configMux.schema)
	if err != nil {
		return nil, err
	}

	categoryHandler, err := categorydelivery.NewCategoryHandler(categoryStorage,
		configMux.addrOrigin, configMux.schema, configMux.portServer,
	)
	if err != nil {
		return nil, err
	}

	productHandler, err := productdelivery.NewProductHandler(productStorage,
		configMux.addrOrigin, configMux.schema, configMux.portServer,
	)
	if err != nil {
		return nil, err
	}

	fileStorage := filerepo.NewFileSystemStorage(configMux.fileServiceDir)
	fileService := fileusecases.NewFileService(fileStorage, urlPrefixPathFS)
	fileHandler := filedelivery.NewFileHandler(fileService, logger,
		configMux.fileServiceDir, configMux.addrOrigin, configMux.schema)

	router.Handle("/api/v1/img/", fileHandler.DocHandlerFileServer())
	router.Handle("/api/v1/img/upload", middleware.Context(ctx, fileHandler.UploadFileHandler))

	router.Handle("/api/v1/signup", middleware.Context(ctx, userHandler.SignUpHandler))
	router.Handle("/api/v1/signin", middleware.Context(ctx, userHandler.SignInHandler))
	router.Handle("/api/v1/logout", middleware.Context(ctx, userHandler.LogOutHandler))

	router.Handle("/api/v1/profile/get", middleware.Context(ctx, userHandler.GetUserHandler))
	router.Handle("/api/v1/profile/update", middleware.Context(ctx, userHandler.PartiallyUpdateUserHandler))

	router.Handle("/api/v1/product/add", middleware.Context(ctx, productHandler.AddProductHandler))
	router.Handle("/api/v1/product/get", middleware.Context(ctx, productHandler.GetProductHandler))
	router.Handle("/api/v1/product/get_list", middleware.Context(ctx, productHandler.GetProductListHandler))
	router.Handle("/api/v1/product/get_list_of_saler",
		middleware.Context(ctx, productHandler.GetListProductOfSalerHandler))
	router.Handle("/api/v1/product/get_list_of_another_saler",
		middleware.Context(ctx, productHandler.GetListProductOfAnotherSalerHandler))
	router.Handle("/api/v1/product/update", middleware.Context(ctx, productHandler.UpdateProductHandler))
	router.Handle("/api/v1/product/close", middleware.Context(ctx, productHandler.CloseProductHandler))
	router.Handle("/api/v1/product/activate", middleware.Context(ctx, productHandler.ActivateProductHandler))
	router.Handle("/api/v1/product/delete", middleware.Context(ctx, productHandler.DeleteProductHandler))

	router.Handle("/api/v1/order/add", middleware.Context(ctx, productHandler.AddOrderHandler))
	router.Handle("/api/v1/order/get_basket", middleware.Context(ctx, productHandler.GetBasketHandler))
	router.Handle("/api/v1/order/update_count", middleware.Context(ctx, productHandler.UpdateOrderCountHandler))
	router.Handle("/api/v1/order/update_status", middleware.Context(ctx, productHandler.UpdateOrderStatusHandler))
	router.Handle("/api/v1/order/buy_full_basket", middleware.Context(ctx, productHandler.BuyFullBasketHandler))
	router.Handle("/api/v1/order/delete", middleware.Context(ctx, productHandler.DeleteOrderHandler))

	router.Handle("/api/v1/category/get_full", middleware.Context(ctx, categoryHandler.GetFullCategories))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router, logger))

	return mux, nil
}

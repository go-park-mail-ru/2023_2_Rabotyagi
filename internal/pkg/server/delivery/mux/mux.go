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
) http.Handler {
	router := http.NewServeMux()

	userHandler := userdelivery.NewUserHandler(userStorage, logger)

	categoryHandler := categorydelivery.NewCategoryHandler(categoryStorage, logger)

	productHandler := productdelivery.NewProductHandler(productStorage, logger)

	fileStorage := filerepo.NewFileSystemStorage(configMux.fileServiceDir)
	fileService := fileusecases.NewFileService(fileStorage, urlPrefixPathFS)
	fileHandler := filedelivery.NewFileHandler(fileService, logger,
		configMux.fileServiceDir, configMux.addrOrigin, configMux.schema)

	router.Handle("/api/v1/img/", fileHandler.DocFileServerHandler(ctx))
	router.Handle("/api/v1/img/upload", middleware.Context(ctx, http.HandlerFunc(fileHandler.UploadFileHandler)))

	router.Handle("/api/v1/signup", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.SignUpHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/signin", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.SignInHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/logout", middleware.Context(ctx, http.HandlerFunc(userHandler.LogOutHandler)))

	router.Handle("/api/v1/profile/get", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.GetUserHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/profile/update", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.PartiallyUpdateUserHandler, configMux.addrOrigin, configMux.schema)))

	router.Handle("/api/v1/product/add", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.AddProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/get", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/get_list", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetProductListHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/get_list_of_saler", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetListProductOfSalerHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/get_list_of_another_saler", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetListProductOfAnotherSalerHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/update", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.UpdateProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/close", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.CloseProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/activate", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.ActivateProductHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/product/delete", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.DeleteProductHandler, configMux.addrOrigin, configMux.schema)))

	router.Handle("/api/v1/order/add", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.AddOrderHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/order/get_basket", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.GetBasketHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/order/update_count", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.UpdateOrderCountHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/order/update_status", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.UpdateOrderStatusHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/order/buy_full_basket", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.BuyFullBasketHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/order/delete", middleware.Context(ctx,
		middleware.SetupCORS(productHandler.DeleteOrderHandler, configMux.addrOrigin, configMux.schema)))

	router.Handle("/api/v1/category/get_full", middleware.Context(ctx,
		middleware.SetupCORS(categoryHandler.GetFullCategories, configMux.addrOrigin, configMux.schema)))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router, logger))

	return mux
}

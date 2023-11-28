package mux

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"net/http"

	categorydelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/delivery"
	citydelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/delivery"
	productdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
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

func NewMux(ctx context.Context, configMux *ConfigMux, userService userdelivery.IUserService,
	productService productdelivery.IProductService, categoryService categorydelivery.ICategoryService,
	cityService citydelivery.ICityService, authGrpcService auth.SessionMangerClient, logger *my_logger.MyLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	authHandler, err := userdelivery.NewAuthHandler(authGrpcService)
	if err != nil {
		return nil, err
	}

	userHandler, err := userdelivery.NewUserHandler(userService)
	if err != nil {
		return nil, err
	}

	categoryHandler, err := categorydelivery.NewCategoryHandler(categoryService)
	if err != nil {
		return nil, err
	}

	cityHandler, err := citydelivery.NewCityHandler(cityService)
	if err != nil {
		return nil, err
	}

	productHandler, err := productdelivery.NewProductHandler(productService)
	if err != nil {
		return nil, err
	}

	router.Handle("/api/v1/signup",
		middleware.SetupCORS(authHandler.SignUpHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/signin",
		middleware.SetupCORS(authHandler.SignInHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/logout", http.HandlerFunc(authHandler.LogOutHandler))

	router.Handle("/api/v1/profile/get",
		middleware.SetupCORS(userHandler.GetUserHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/profile/update",
		middleware.SetupCORS(userHandler.PartiallyUpdateUserHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/api/v1/product/add",
		middleware.SetupCORS(productHandler.AddProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/get",
		middleware.SetupCORS(productHandler.GetProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/get_list",
		middleware.SetupCORS(productHandler.GetProductListHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/get_list_of_saler",
		middleware.SetupCORS(productHandler.GetListProductOfSalerHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/get_list_of_another_saler",
		middleware.SetupCORS(productHandler.GetListProductOfAnotherSalerHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/update",
		middleware.SetupCORS(productHandler.UpdateProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/close",
		middleware.SetupCORS(productHandler.CloseProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/activate",
		middleware.SetupCORS(productHandler.ActivateProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/delete",
		middleware.SetupCORS(productHandler.DeleteProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/search",
		middleware.SetupCORS(productHandler.SearchProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/get_search_feed",
		middleware.SetupCORS(productHandler.GetSearchProductFeedHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/api/v1/profile/favourites",
		middleware.SetupCORS(productHandler.GetFavouritesHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/add-to-fav",
		middleware.SetupCORS(productHandler.AddToFavouritesHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/product/remove-from-fav",
		middleware.SetupCORS(productHandler.DeleteFromFavouritesHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/api/v1/order/add",
		middleware.SetupCORS(productHandler.AddOrderHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/order/get_basket",
		middleware.SetupCORS(productHandler.GetBasketHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/order/update_count",
		middleware.SetupCORS(productHandler.UpdateOrderCountHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/order/update_status",
		middleware.SetupCORS(productHandler.UpdateOrderStatusHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/order/buy_full_basket",
		middleware.SetupCORS(productHandler.BuyFullBasketHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/order/delete",
		middleware.SetupCORS(productHandler.DeleteOrderHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/api/v1/category/get_full",
		middleware.SetupCORS(categoryHandler.GetFullCategories, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/category/search",
		middleware.SetupCORS(categoryHandler.SearchCategoryHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/api/v1/city/get_full",
		middleware.SetupCORS(cityHandler.GetFullCitiesHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/api/v1/city/search",
		middleware.SetupCORS(cityHandler.SearchCityHandler, configMux.addrOrigin, configMux.schema))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.Context(ctx,
		middleware.AddReqID(middleware.AccessLogMiddleware(router, logger))), logger))

	return mux, nil
}

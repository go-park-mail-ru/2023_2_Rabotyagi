package mux

import (
	"context"
	"net/http"

	categorydelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/delivery"
	citydelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/delivery"
	productdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ConfigMux struct {
	addrOrigin      string
	schema          string
	portServer      string
	mainServiceName string
}

func NewConfigMux(addrOrigin string, schema string, portServer string, mainServiceName string) *ConfigMux {
	return &ConfigMux{
		addrOrigin:      addrOrigin,
		schema:          schema,
		portServer:      portServer,
		mainServiceName: mainServiceName,
	}
}

//nolint:funlen
func NewMux(ctx context.Context, configMux *ConfigMux, userService userdelivery.IUserService,
	productService productdelivery.IProductService, categoryService categorydelivery.ICategoryService,
	cityService citydelivery.ICityService, authGrpcService auth.SessionMangerClient,
	logger *mylogger.MyLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	authHandler, err := userdelivery.NewAuthHandler(authGrpcService)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	userHandler, err := userdelivery.NewProfileHandler(userService, authGrpcService)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	categoryHandler, err := categorydelivery.NewCategoryHandler(categoryService)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	cityHandler, err := citydelivery.NewCityHandler(cityService)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	productHandler, err := productdelivery.NewProductHandler(productService, authGrpcService)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	router.Handle("/signup",
		middleware.SetupCORS(authHandler.SignUpHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/signin",
		middleware.SetupCORS(authHandler.SignInHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/logout", http.HandlerFunc(authHandler.LogOutHandler))

	router.Handle("/profile/get",
		middleware.SetupCORS(userHandler.GetUserHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/profile/update",
		middleware.SetupCORS(userHandler.PartiallyUpdateUserHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/product/add",
		middleware.SetupCORS(productHandler.AddProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/get",
		middleware.SetupCORS(productHandler.GetProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/get_list",
		middleware.SetupCORS(productHandler.GetProductListHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/get_list_of_saler",
		middleware.SetupCORS(productHandler.GetListProductOfSalerHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/get_list_of_another_saler",
		middleware.SetupCORS(productHandler.GetListProductOfAnotherSalerHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/update",
		middleware.SetupCORS(productHandler.UpdateProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/close",
		middleware.SetupCORS(productHandler.CloseProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/activate",
		middleware.SetupCORS(productHandler.ActivateProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/delete",
		middleware.SetupCORS(productHandler.DeleteProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/search",
		middleware.SetupCORS(productHandler.SearchProductHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/get_search_feed",
		middleware.SetupCORS(productHandler.GetSearchProductFeedHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/profile/favourites",
		middleware.SetupCORS(productHandler.GetFavouritesHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/add-to-fav",
		middleware.SetupCORS(productHandler.AddToFavouritesHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/product/remove-from-fav",
		middleware.SetupCORS(productHandler.DeleteFromFavouritesHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/premium/add",
		middleware.SetupCORS(productHandler.AddPremiumHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/order/add",
		middleware.SetupCORS(productHandler.AddOrderHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/get_basket",
		middleware.SetupCORS(productHandler.GetBasketHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/get_not_in_basket",
		middleware.SetupCORS(productHandler.GetBasketHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/sold",
		middleware.SetupCORS(productHandler.GetBasketHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/update_count",
		middleware.SetupCORS(productHandler.UpdateOrderCountHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/update_status",
		middleware.SetupCORS(productHandler.UpdateOrderStatusHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/buy_full_basket",
		middleware.SetupCORS(productHandler.BuyFullBasketHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/order/delete",
		middleware.SetupCORS(productHandler.DeleteOrderHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/category/get_full",
		middleware.SetupCORS(categoryHandler.GetFullCategories, configMux.addrOrigin, configMux.schema))
	router.Handle("/category/search",
		middleware.SetupCORS(categoryHandler.SearchCategoryHandler, configMux.addrOrigin, configMux.schema))

	router.Handle("/city/get_full",
		middleware.SetupCORS(cityHandler.GetFullCitiesHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/city/search",
		middleware.SetupCORS(cityHandler.SearchCityHandler, configMux.addrOrigin, configMux.schema))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/healthcheck", delivery.HealthCheckHandler)

	metricsManager := metrics.NewMetricManagerHTTP(configMux.mainServiceName)
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.Context(ctx,
		middleware.AddReqID(
			middleware.AccessLogMiddleware(
				middleware.AddAPIName(router, middleware.APINameV1),
				logger, metricsManager))),
		logger))

	return mux, nil
}

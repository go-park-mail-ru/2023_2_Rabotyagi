package delivery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ IProductService = (*usecases.ProductService)(nil)

type IProductService interface {
	AddProduct(ctx context.Context, r io.Reader, userID uint64) (productID uint64, err error)
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error)
	GetProductsList(ctx context.Context,
		lastProductID uint64, count uint64, userID uint64) ([]*models.ProductInFeed, error)
	GetProductsOfSaler(ctx context.Context, lastProductID uint64,
		count uint64, userID uint64, isMy bool) ([]*models.ProductInFeed, error)
	UpdateProduct(ctx context.Context, r io.Reader, isPartialUpdate bool, productID uint64, userAuthID uint64) error
	CloseProduct(ctx context.Context, productID uint64, userID uint64) error
	ActivateProduct(ctx context.Context, productID uint64, userID uint64) error
	DeleteProduct(ctx context.Context, productID uint64, userID uint64) error
	SearchProduct(ctx context.Context, searchInput string) ([]string, error)
	GetSearchProductFeed(ctx context.Context,
		searchInput string, lastNumber uint64, limit uint64, userID uint64,
	) ([]*models.ProductInFeed, error)
	IBasketService
	IFavouriteService
}

type ProductHandler struct {
	sessionManagerClient auth.SessionMangerClient
	service              IProductService
	logger               *my_logger.MyLogger
}

func NewProductHandler(productService IProductService,
	sessionManagerClient auth.SessionMangerClient,
) (*ProductHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductHandler{
		service:              productService,
		logger:               logger,
		sessionManagerClient: sessionManagerClient,
	}, nil
}

// AddProductHandler godoc
//
//	@Summary    add product
//	@Description  add product by data
//	@Tags product
//
//	@Accept      json
//	@Produce    json
//	@Param      product  body models.PreProduct true  "product data for adding"
//	@Success    200  {object} responses.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /product/add [post]
func (p *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	productID, err := p.service.AddProduct(ctx, r.Body, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, responses.NewResponseIDRedirect(productID))
	logger.Infof("in AddProductHandler: added product id= %+v", productID)
}

// GetProductHandler godoc
//
//	@Summary    get product
//	@Description  get product by id
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Success    200  {object} ProductResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/get [get]
func (p *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		if errors.Is(err, responses.ErrCookieNotPresented) {
			userID = 0
		} else {
			responses.HandleErr(w, logger, err)

			return
		}
	}

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	product, err := p.service.GetProduct(ctx, productID, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProductResponse(product))
	logger.Infof("in GetProductHandler: get product: %+v", product)
}

// GetProductListHandler godoc
//
//	@Summary    get products list
//	@Description  get products by count and last_id return old products
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count products"
//	@Param      last_id  query uint64 true  "last product id"
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/get_list [get]
func (p *ProductHandler) GetProductListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	count, err := utils.ParseUint64FromRequest(r, "count")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	lastID, err := utils.ParseUint64FromRequest(r, "last_id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		if errors.Is(err, responses.ErrCookieNotPresented) {
			userID = 0
		} else {
			responses.HandleErr(w, logger, err)

			return
		}
	}

	products, err := p.service.GetProductsList(ctx, lastID, count, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProductListResponse(products))
	logger.Infof("in GetProductListHandler: get product list: %+v", products)
}

// GetListProductOfSalerHandler godoc
//
//	@Summary     get list of products for saler
//	@Description  get list of products for saler using user id from cookies\jwt
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count products"
//	@Param      last_id  query uint64 true  "last product id "
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/get_list_of_saler [get]
func (p *ProductHandler) GetListProductOfSalerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	count, err := utils.ParseUint64FromRequest(r, "count")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	lastID, err := utils.ParseUint64FromRequest(r, "last_id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	products, err := p.service.GetProductsOfSaler(ctx, lastID, count, userID, true)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProductListResponse(products))
	logger.Infof("in GetListProductOfSalerHandler: get product list: %+v", products)
}

// GetListProductOfAnotherSalerHandler godoc
//
//	@Summary     get list of products for another saler
//	@Description  get list of products for another saler using saler id, count and last product id from query
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      saler_id  query uint64 true  "saler id"
//	@Param      count  query uint64 true  "count products"
//	@Param      last_id  query uint64 true  "last product id "
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/get_list_of_another_saler [get]
func (p *ProductHandler) GetListProductOfAnotherSalerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	count, err := utils.ParseUint64FromRequest(r, "count")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	lastID, err := utils.ParseUint64FromRequest(r, "last_id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	salerID, err := utils.ParseUint64FromRequest(r, "saler_id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	products, err := p.service.GetProductsOfSaler(ctx, lastID, count, salerID, false)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProductListResponse(products))
	logger.Infof("in GetListProductOfAnotherSalerHandler: get product list: %+v", products)
}

// UpdateProductHandler godoc
//
//	@Summary    update product
//	@Description  update product by id
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Param      preProduct  body models.PreProduct false  "полностью опционален"
//	@Success    200  {object} responses.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /product/update [patch]
//	@Router      /product/update [put]
func (p *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	if r.Method == http.MethodPatch {
		err = p.service.UpdateProduct(ctx, r.Body, true, productID, userID)
	} else {
		err = p.service.UpdateProduct(ctx, r.Body, false, productID, userID)
	}

	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, responses.NewResponseIDRedirect(productID))
	logger.Infof("in UpdateProductHandler: updated product with id = %+v", productID)
}

// CloseProductHandler godoc
//
//	@Summary     close product
//	@Description  close product for saler using user id from cookies\jwt.
//	@Description  This does product not active.
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/close [patch]
func (p *ProductHandler) CloseProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	err = p.service.CloseProduct(ctx, productID, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulCloseProduct))
	logger.Infof("in CloseProductHandler: close product id=%d", productID)
}

// ActivateProductHandler godoc
//
//	@Summary     activate product
//	@Description  activate product for saler using user id from cookies\jwt.
//	@Description  This does product active.
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/activate [patch]
func (p *ProductHandler) ActivateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	err = p.service.ActivateProduct(ctx, productID, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulActivateProduct))
	logger.Infof("in ActivateProductHandler: activated product id=%d", productID)
}

// DeleteProductHandler godoc
//
//	@Summary     delete product
//	@Description  delete product for saler using user id from cookies\jwt.
//	@Description  This totally removed product. Recovery will be impossible
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badContent(4400)
//	@Router      /product/delete [delete]
func (p *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	err = p.service.DeleteProduct(ctx, productID, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulDeleteProduct))
	logger.Infof("in DeleteProductHandler: delete product id=%d", productID)
}

// SearchProductHandler godoc
//
//	@Summary    search products
//	@Description  search top 5 common named/descripted products
//	@Tags product
//	@Produce    json
//	@Param      searched  query string true  "searched string"
//	@Success    200  {object} ProductInSearchListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/search [get]
func (p *ProductHandler) SearchProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	searchInput := utils.ParseStringFromRequest(r, "searched")

	products, err := p.service.SearchProduct(ctx, searchInput)
	if err != nil {
		responses.SendResponse(w, logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	responses.SendResponse(w, logger, NewProductInSearchListResponse(products))
	logger.Infof("in SearchProductHandler: search products: %+v\n", products)
}

// GetSearchProductFeedHandler godoc
//
//	@Summary    get products search feed
//	@Description  get products feed after search
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count products"
//	@Param      offset  query uint64 true  "last product id"
//	@Param      searched  query string true  "searched string"
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error" Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/get_search_feed [get]
func (p *ProductHandler) GetSearchProductFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	count, err := utils.ParseUint64FromRequest(r, "count")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	offset, err := utils.ParseUint64FromRequest(r, "offset")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	searchInput := utils.ParseStringFromRequest(r, "searched")

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		if errors.Is(err, responses.ErrCookieNotPresented) {
			userID = 0
		} else {
			responses.HandleErr(w, logger, err)

			return
		}
	}

	products, err := p.service.GetSearchProductFeed(ctx, searchInput, offset, count, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProductListResponse(products))
	logger.Infof("in GetProductListHandler: get product list: %+v", products)
}

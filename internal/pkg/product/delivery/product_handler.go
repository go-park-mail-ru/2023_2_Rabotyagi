package delivery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"

	"go.uber.org/zap"
)

var _ IProductService = (*usecases.ProductService)(nil)

type IProductService interface {
	AddProduct(ctx context.Context, r io.Reader) (productID uint64, err error)
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error)
	GetProductsList(ctx context.Context,
		lastProductID uint64, count uint64, userID uint64) ([]*models.ProductInFeed, error)
	GetProductsOfSaler(ctx context.Context, lastProductID uint64,
		count uint64, userID uint64, isMy bool) ([]*models.ProductInFeed, error)
	UpdateProduct(ctx context.Context, r io.Reader, isPartialUpdate bool, productID uint64, userAuthID uint64) error
	CloseProduct(ctx context.Context, productID uint64, userID uint64) error
	ActivateProduct(ctx context.Context, productID uint64, userID uint64) error
	DeleteProduct(ctx context.Context, productID uint64, userID uint64) error
	IBasketService
	IFavouriteService
}

type ProductHandler struct {
	service IProductService
	logger  *zap.SugaredLogger
}

func NewProductHandler(productService IProductService) (*ProductHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductHandler{
		service: productService,
		logger:  logger,
	}, nil
}

// AddProductHandler godoc
//
//	@Summary    add product
//	@Description  add product by data
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags product
//
//	@Accept      json
//	@Produce    json
//	@Param      product  body models.PreProduct true  "product data for adding"
//	@Success    200  {object} delivery.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/add [post]
func (p *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	productID, err := p.service.AddProduct(ctx, r.Body)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, delivery.NewResponseID(productID))
	p.logger.Infof("in AddProductHandler: added product id= %+v", productID)
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get [get]
func (p *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		if errors.Is(err, delivery.ErrCookieNotPresented) {
			userID = 0
		} else {
			delivery.HandleErr(w, p.logger, err)

			return
		}
	}

	productID, err := parseIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	product, err := p.service.GetProduct(ctx, productID, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductResponse(delivery.StatusResponseSuccessful, product))
	p.logger.Infof("in GetProductHandler: get product: %+v", product)
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get_list [get]
func (p *ProductHandler) GetProductListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	count, lastID, err := parseCountAndLastIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		if errors.Is(err, delivery.ErrCookieNotPresented) {
			userID = 0
		} else {
			delivery.HandleErr(w, p.logger, err)

			return
		}
	}

	products, err := p.service.GetProductsList(ctx, lastID, count, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetProductListHandler: get product list: %+v", products)
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get_list_of_saler [get]
func (p *ProductHandler) GetListProductOfSalerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	count, lastID, err := parseCountAndLastIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	products, err := p.service.GetProductsOfSaler(ctx, lastID, count, userID, true)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetListProductOfSalerHandler: get product list: %+v", products)
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get_list_of_another_saler [get]
func (p *ProductHandler) GetListProductOfAnotherSalerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	salerID, count, lastID, err := parseSalerIDCountLastIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	products, err := p.service.GetProductsOfSaler(ctx, lastID, count, salerID, false)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetListProductOfAnotherSalerHandler: get product list: %+v", products)
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
//	@Success    200  {object} delivery.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/update [patch]
//	@Router      /product/update [put]
func (p *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	productID, err := parseIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	if r.Method == http.MethodPatch {
		err = p.service.UpdateProduct(ctx, r.Body, true, productID, userID)
	} else {
		err = p.service.UpdateProduct(ctx, r.Body, false, productID, userID)
	}

	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, delivery.NewResponseID(productID))
	p.logger.Infof("in UpdateProductHandler: updated product with id = %+v", productID)
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
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/close [patch]
func (p *ProductHandler) CloseProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	productID, err := parseIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.CloseProduct(ctx, productID, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulCloseProduct))
	p.logger.Infof("in CloseProductHandler: close product id=%d", productID)
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
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/activate [patch]
func (p *ProductHandler) ActivateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	productID, err := parseIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.ActivateProduct(ctx, productID, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulCloseProduct))
	p.logger.Infof("in ActivateProductHandler: activated product id=%d", productID)
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
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/delete [delete]
func (p *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	productID, err := parseIDFromRequest(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.DeleteProduct(ctx, productID, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulDeleteProduct))
	p.logger.Infof("in DeleteProductHandler: delete product id=%d", productID)
}

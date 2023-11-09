package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"

	"go.uber.org/zap"
)

type ProductHandler struct {
	storage    usecases.IProductStorage
	addrOrigin string
	schema     string
	portServer string
	logger     *zap.SugaredLogger
}

func NewProductHandler(storage usecases.IProductStorage,
	addrOrigin string, schema string, portServer string, logger *zap.SugaredLogger,
) *ProductHandler {
	return &ProductHandler{
		storage:    storage,
		addrOrigin: addrOrigin,
		schema:     schema,
		portServer: portServer,
		logger:     logger,
	}
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
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	preProduct, err := usecases.ValidatePreProduct(p.logger, r.Body)
	if err != nil {
		p.logger.Errorf("in AddProductHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	productID, err := p.storage.AddProduct(ctx, preProduct)
	if err != nil {
		p.logger.Errorf("in AddProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, p.logger, delivery.NewResponseID(productID))
	p.logger.Infof("in AddProductHandler: added product: %+v", preProduct)
}

// GetProductHandler godoc
//
//	@Summary    get product
//	@Description  get product by id
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "product id"
//	@Success    200  {object} ProductResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get/{id} [get]
func (p *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	productIDStr := delivery.GetPathParam(r.URL.Path)
	userID := delivery.GetUserIDFromCookie(r, p.logger)

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		p.logger.Errorf("in GetProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s product id == %s But shoud be integer", delivery.ErrBadRequest, productIDStr)))

		return
	}

	product, err := p.storage.GetProduct(ctx, productID, userID)
	if err != nil {
		p.logger.Errorf("in GetProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrProductNotExist))

		return
	}

	product.Sanitize()

	delivery.SendOkResponse(w, p.logger, NewProductResponse(delivery.StatusResponseSuccessful, product))
	p.logger.Infof("in GetProductHandler: get product: %+v", product)
}

// GetProductListHandler godoc
//
//	@Summary    get products list
//	@Description  get products by count and last_id return new products
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count products"
//	@Param      last_id  query uint64 true  "last product id "
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get_list [get]
func (p *ProductHandler) GetProductListHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	count, lastID, err := parseCountAndLastIDFromRequest(r, p.logger)
	if err != nil {
		p.logger.Errorf("in GetListProductOfSalerHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r, p.logger)

	products, err := p.storage.GetNewProducts(ctx, lastID, count, userID)
	if err != nil {
		p.logger.Errorf("in GetListProductOfSalerHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	for _, product := range products {
		product.Sanitize()
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetProductListHandler: get product list: %+v", products)
}

// UpdateProductHandler godoc
//
//	@Summary    update product
//	@Description  update product by id
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      product_id  path uint64 true  "id of product"
//	@Param      preProduct  body models.PreProduct true  "product data for updating"
//	@Success    200  {object} delivery.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/update [patch]
//	@Router      /product/update [put]
func (p *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	productIDStr := delivery.GetPathParam(r.URL.Path)

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		p.logger.Errorf("in UpdateProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s product id == %s But shoud be integer", delivery.ErrBadRequest, productIDStr)))

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r, p.logger)

	var preProduct *models.PreProduct

	if r.Method == http.MethodPatch {
		preProduct, err = usecases.ValidatePartOfPreProduct(p.logger, r.Body)
		if err != nil {
			p.logger.Errorf("in UpdateProductHandler: %+v\n", err)
			delivery.HandleErr(w, p.logger, err)

			return
		}
	} else {
		preProduct, err = usecases.ValidatePreProduct(p.logger, r.Body)
		if err != nil {
			p.logger.Errorf("in UpdateProductHandler: %+v\n", err)
			delivery.HandleErr(w, p.logger, err)

			return
		}
	}

	if preProduct.SalerID != userID {
		delivery.SendErrResponse(w, p.logger, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrUserPermissionsChange))

		return
	}

	updateFieldsMap := utils.StructToMap(preProduct)

	err = p.storage.UpdateProduct(ctx, productID, updateFieldsMap)
	if err != nil {
		p.logger.Errorf("in UpdateProductHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, delivery.NewResponseID(productID))
	p.logger.Infof("in UpdateProductHandler: updated product with id = %+v", productID)
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
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	count, lastID, err := parseCountAndLastIDFromRequest(r, p.logger)
	if err != nil {
		p.logger.Errorf("in GetListProductOfSalerHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r, p.logger)

	products, err := p.storage.GetProductsOfSaler(ctx, lastID, count, userID)
	if err != nil {
		p.logger.Errorf("in GetListProductOfSalerHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	for _, product := range products {
		product.Sanitize()
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetListProductOfSalerHandler: get product list: %+v", products)
}

// CloseProductHandler godoc
//
//	@Summary     close product
//	@Description  close product for saler using user id from cookies\jwt.
//	@Description  This does product not active.
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      productID  path uint64 true  "product id"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/close/ [patch]
func (p *ProductHandler) CloseProductHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r, p.logger)
	productIDStr := delivery.GetPathParam(r.URL.String())

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		p.logger.Errorf("in CloseProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrWrongProductID.Error()))

		return
	}

	err = p.storage.CloseProduct(ctx, productID, userID)
	if err != nil {
		p.logger.Errorf("in CloseProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulCloseProduct))
	p.logger.Infof("in CloseProductHandler: close product id=%d", productID)
}

// DeleteProductHandler godoc
//
//	@Summary     delete product
//	@Description  delete product for saler using user id from cookies\jwt.
//	@Description  This totally removed product. Recovery will be impossible
//	@Tags product
//	@Accept      json
//	@Produce    json
//	@Param      productID  path uint64 true  "product id"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/delete/ [delete]
func (p *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r, p.logger)
	productIDStr := delivery.GetPathParam(r.URL.String())

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		p.logger.Errorf("in DeleteProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrWrongProductID.Error()))

		return
	}

	err = p.storage.DeleteProduct(ctx, productID, userID)
	if err != nil {
		p.logger.Errorf("in DeleteProductHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulDeleteProduct))
	p.logger.Infof("in DeleteProductHandler: delete product id=%d", productID)
}

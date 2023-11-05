package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
)

type ProductHandler struct {
	storage    usecases.IProductStorage
	addrOrigin string
	schema     string
	portServer string
}

func NewProductHandler(storage usecases.IProductStorage,
	addrOrigin string, schema string, portServer string,
) *ProductHandler {
	return &ProductHandler{
		storage:    storage,
		addrOrigin: addrOrigin,
		schema:     schema,
		portServer: portServer,
	}
}

func (p *ProductHandler) createURLToProductFromID(productID uint64) string {
	return fmt.Sprintf("%s%s:%s/api/v1/product/get/%d", p.schema, p.addrOrigin, p.portServer, productID)
}

// AddProductHandler godoc
//
//	@Summary    add product
//	@Description  add product by data
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//
//	@Accept      json
//	@Produce    json
//	@Param      product  body models.PreProduct true  "product data for adding"
//	@Success    200  {object} delivery.ResponseRedirect
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/add [post]
func (p *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	preProduct, err := usecases.ValidatePreProduct(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in AddProductHandler:", err)

		return
	}

	productID, err := p.storage.AddProduct(ctx, preProduct)
	if err != nil {
		log.Printf("in AddProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, delivery.NewResponseRedirect(p.createURLToProductFromID(productID)))
	log.Printf("added product: %+v", preProduct)
}

// GetProductHandler godoc
//
//	@Summary    get product
//	@Description  get product by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "product id"
//	@Success    200  {object} ProductResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get/{id} [get]
func (p *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	ctx := r.Context()
	productIDStr := utils.GetPathParam(r.URL.Path)

	userID := usecases.GetUserIDFromCookie(r)

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		log.Printf("in GetProductHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s product id == %s But shoud be integer", delivery.ErrBadRequest, productIDStr)))

		return
	}

	product, err := p.storage.GetProduct(ctx, productID, userID)
	if err != nil {
		log.Printf("in GetProductHandler: product with this id is not exists %+v\n", productID)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrProductNotExist))

		return
	}

	delivery.SendOkResponse(w, NewProductResponse(delivery.StatusResponseSuccessful, product))
	log.Printf("in GetProductHandler: get product: %+v", product)
}

// GetProductListHandler godoc
//
//	@Summary    get product
//	@Description  get product by count
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
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	countStr := r.URL.Query().Get("count")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		log.Printf("in GetProductListHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s count products == %s But shoud be integer", delivery.ErrBadRequest, countStr)))

		return
	}

	lastIDStr := r.URL.Query().Get("last_id")

	lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
	if err != nil {
		log.Printf("in GetProductListHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s last_id products == %s But shoud be integer", delivery.ErrBadRequest, lastIDStr)))

		return
	}

	ctx := r.Context()

	userID := usecases.GetUserIDFromCookie(r)

	products, err := p.storage.GetNewProducts(ctx, lastID, count, userID)
	if err != nil {
		log.Printf("in GetProductListHandler %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	log.Printf("in GetProductListHandler: get product list: %+v", products)
}

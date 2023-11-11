package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"go.uber.org/zap"
	"net/http"
)

type CategoryHandler struct {
	storage    usecases.ICategoryStorage
	addrOrigin string
	schema     string
	portServer string
	logger     *zap.SugaredLogger
}

func NewCategoryHandler(storage usecases.ICategoryStorage,
	addrOrigin string, schema string, portServer string, logger *zap.SugaredLogger,
) *CategoryHandler {
	return &CategoryHandler{
		storage:    storage,
		addrOrigin: addrOrigin,
		schema:     schema,
		portServer: portServer,
		logger:     logger,
	}
}

// GetFullCategories godoc
//
//	@Summary    get all categories
//	@Description  get all categories
//	@Tags category
//	@Produce    json
//	@Success    200  {object} CategoryListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /category/get_full [get]
func (c *CategoryHandler) GetFullCategories(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, c.addrOrigin, c.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	categories, err := c.storage.GetFullCategories(ctx)
	if err != nil {
		c.logger.Errorf("in GetFullCategories %+v\n", err)
		delivery.SendErrResponse(w, c.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	for _, order := range categories {
		order.Sanitize()
	}

	delivery.SendOkResponse(w, c.logger, NewCategoryListResponse(delivery.StatusResponseSuccessful, categories))
	c.logger.Infof("in GetBasketHandler: get basket of orders: %+v\n", categories)
}

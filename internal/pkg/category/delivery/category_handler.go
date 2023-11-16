package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/category/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"

	"go.uber.org/zap"
)

var _ ICategoryService = (*usecases.CategoryService)(nil)

type ICategoryService interface {
	GetFullCategories(ctx context.Context) ([]*models.Category, error)
}

type CategoryHandler struct {
	service ICategoryService
	logger  *zap.SugaredLogger
}

func NewCategoryHandler(service ICategoryService) (*CategoryHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &CategoryHandler{
		service: service,
		logger:  logger,
	}, nil
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
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	categories, err := c.service.GetFullCategories(ctx)
	if err != nil {
		delivery.SendErrResponse(w, c.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, c.logger, NewCategoryListResponse(delivery.StatusResponseSuccessful, categories))
	c.logger.Infof("in GetBasketHandler: get basket of orders: %+v\n", categories)
}

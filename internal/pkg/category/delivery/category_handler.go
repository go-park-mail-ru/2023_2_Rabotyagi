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
	SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error)
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
	c.logger.Infof("in GetFullCategories: get all categories: %+v\n", categories)
}

// SearchCategoryHandler godoc
//
//	@Summary    search category
//	@Description  search top 5 common named categories
//	@Tags category
//	@Produce    json
//	@Param      searched  query string true  "searched string"
//	@Success    200  {object} CategoryListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /category/search [get]
func (c *CategoryHandler) SearchCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	searchInput := r.URL.Query().Get("searched")

	categories, err := c.service.SearchCategory(ctx, searchInput)
	if err != nil {
		delivery.SendErrResponse(w, c.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, c.logger, NewCategoryListResponse(delivery.StatusResponseSuccessful, categories))
	c.logger.Infof("in SearchCategoryHandler: search category: %+v\n", categories)
}

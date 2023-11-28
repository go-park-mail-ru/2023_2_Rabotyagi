package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"net/http"

	"go.uber.org/zap"
)

var _ ICityService = (*usecases.CityService)(nil)

type ICityService interface {
	GetFullCities(ctx context.Context) ([]*models.City, error)
	SearchCity(ctx context.Context, searchInput string) ([]*models.City, error)
}

type CityHandler struct {
	service ICityService
	logger  *zap.SugaredLogger
}

func NewCityHandler(service ICityService) (*CityHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &CityHandler{
		service: service,
		logger:  logger,
	}, nil
}

// GetFullCitiesHandler godoc
//
//	@Summary    get all cities
//	@Description  get all cities
//	@Tags city
//	@Produce    json
//	@Success    200  {object} CityListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /city/get_full [get]
func (c *CityHandler) GetFullCitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	cities, err := c.service.GetFullCities(ctx)
	if err != nil {
		responses.SendResponse(w, c.logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	responses.SendResponse(w, c.logger, NewCityListResponse(cities))
	c.logger.Infof("in GetFullCities: get all cities: %+v\n", cities)
}

// SearchCityHandler godoc
//
//	@Summary    search city
//	@Description  search top 5 common named cities
//	@Tags city
//	@Produce    json
//	@Param      searched  query string true  "searched string"
//	@Success    200  {object} CityListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /city/search [get]
func (c *CityHandler) SearchCityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	searchInput := utils.ParseStringFromRequest(r, "searched")

	cities, err := c.service.SearchCity(ctx, searchInput)
	if err != nil {
		responses.SendResponse(w, c.logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	responses.SendResponse(w, c.logger, NewCityListResponse(cities))
	c.logger.Infof("in SearchCityHandler: search city: %+v\n", cities)
}

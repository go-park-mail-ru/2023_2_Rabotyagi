package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ ICityService = (*usecases.CityService)(nil)

type ICityService interface {
	GetFullCities(ctx context.Context) ([]*models.City, error)
	SearchCity(ctx context.Context, searchInput string) ([]*models.City, error)
}

type CityHandler struct {
	service ICityService
	logger  *mylogger.MyLogger
}

func NewCityHandler(service ICityService) (*CityHandler, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
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
//	@Failure    222  {object} responses.ErrorResponse "Error"
//	@Router      /city/get_full [get]
func (c *CityHandler) GetFullCitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := c.logger.LogReqID(ctx)

	cities, err := c.service.GetFullCities(ctx)
	if err != nil {
		responses.SendResponse(w, logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	responses.SendResponse(w, logger, NewCityListResponse(cities))
	logger.Infof("in GetFullCities: get all cities: %+v\n", cities)
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
//	@Failure    222  {object} responses.ErrorResponse "Error"
//	@Router      /city/search [get]
func (c *CityHandler) SearchCityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := c.logger.LogReqID(ctx)

	searchInput := utils.ParseStringFromRequest(r, "searched")

	cities, err := c.service.SearchCity(ctx, searchInput)
	if err != nil {
		responses.SendResponse(w, logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	responses.SendResponse(w, logger, NewCityListResponse(cities))
	logger.Infof("in SearchCityHandler: search city: %+v\n", cities)
}

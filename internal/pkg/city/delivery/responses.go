package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/statuses"
)

type CityListResponse struct {
	Status int            `json:"status"`
	Body   []*models.City `json:"body"`
}

func NewCityListResponse(body []*models.City) *CityListResponse {
	return &CityListResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

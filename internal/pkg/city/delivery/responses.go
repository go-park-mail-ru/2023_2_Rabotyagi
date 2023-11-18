package delivery

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"

type CityListResponse struct {
	Status int            `json:"status"`
	Body   []*models.City `json:"body"`
}

func NewCityListResponse(status int, body []*models.City) *CityListResponse {
	return &CityListResponse{
		Status: status,
		Body:   body,
	}
}

package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

//easyjson:json
type CategoryListResponse struct {
	Status int                `json:"status"`
	Body   []*models.Category `json:"body"`
}

func NewCategoryListResponse(body []*models.Category) *CategoryListResponse {
	return &CategoryListResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

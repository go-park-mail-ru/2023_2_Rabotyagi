package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery/statuses"
)

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

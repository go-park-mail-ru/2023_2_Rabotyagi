package delivery

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"

type CategoryListResponse struct {
	Status int                `json:"status"`
	Body   []*models.Category `json:"body"`
}

func NewCategoryListResponse(status int, body []*models.Category) *CategoryListResponse {
	return &CategoryListResponse{
		Status: status,
		Body:   body,
	}
}

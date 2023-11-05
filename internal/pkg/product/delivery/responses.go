package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

const (
	ResponseSuccessfulAddProduct = "Successful add product"

	ErrProductNotExist      = "Product not exists"
	ErrNoSuchCountOfProduct = "not enough products in storage"
)

type ProductResponse struct {
	Status int             `json:"status"`
	Body   *models.Product `json:"body"`
}

func NewProductResponse(status int, body *models.Product) *ProductResponse {
	return &ProductResponse{
		Status: status,
		Body:   body,
	}
}

type ProductListResponse struct {
	Status int                     `json:"status"`
	Body   []*models.ProductInFeed `json:"body"`
}

func NewProductListResponse(status int, body []*models.ProductInFeed) *ProductListResponse {
	return &ProductListResponse{
		Status: status,
		Body:   body,
	}
}

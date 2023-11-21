package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

const (
	ResponseSuccessfulUpdateCountOrder  = "Успешно изменено количество заказа"
	ResponseSuccessfulUpdateStatusOrder = "Успешно изменен статус заказа"
	ResponseSuccessfulBuyFullBasket     = "Успешная покупка всего из корзины"
	ResponseSuccessfulCloseProduct      = "Объявление успешно закрыто"
	ResponseSuccessfulDeleteProduct     = "Объявление успешно удалено"
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

type ProductInSearchListResponse struct {
	Status int                       `json:"status"`
	Body   []*models.ProductInSearch `json:"body"`
}

func NewProductInSearchListResponse(status int, body []*models.ProductInSearch) *ProductInSearchListResponse {
	return &ProductInSearchListResponse{
		Status: status,
		Body:   body,
	}
}

type OrderResponse struct {
	Status int                   `json:"status"`
	Body   *models.OrderInBasket `json:"body"`
}

func NewOrderResponse(status int, body *models.OrderInBasket) *OrderResponse {
	return &OrderResponse{
		Status: status,
		Body:   body,
	}
}

type OrderListResponse struct {
	Status int                     `json:"status"`
	Body   []*models.OrderInBasket `json:"body"`
}

func NewOrderListResponse(status int, body []*models.OrderInBasket) *OrderListResponse {
	return &OrderListResponse{
		Status: status,
		Body:   body,
	}
}

package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/statuses"
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

func NewProductResponse(body *models.Product) *ProductResponse {
	return &ProductResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

type ProductListResponse struct {
	Status int                     `json:"status"`
	Body   []*models.ProductInFeed `json:"body"`
}

func NewProductListResponse(body []*models.ProductInFeed) *ProductListResponse {
	return &ProductListResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

type OrderResponse struct {
	Status int                   `json:"status"`
	Body   *models.OrderInBasket `json:"body"`
}

func NewOrderResponse(body *models.OrderInBasket) *OrderResponse {
	return &OrderResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

type OrderListResponse struct {
	Status int                     `json:"status"`
	Body   []*models.OrderInBasket `json:"body"`
}

func NewOrderListResponse(body []*models.OrderInBasket) *OrderListResponse {
	return &OrderListResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

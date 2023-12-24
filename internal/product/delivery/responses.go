package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

const (
	ResponseSuccessfulUpdateCountOrder  = "Успешно изменено количество заказа"
	ResponseSuccessfulUpdateStatusOrder = "Успешно изменен статус заказа"
	ResponseSuccessfulBuyFullBasket     = "Успешная покупка всего из корзины"
	ResponseSuccessfulCloseProduct      = "Объявление успешно закрыто"
	ResponseSuccessfulDeleteProduct     = "Объявление успешно удалено"
	ResponseSuccessfulActivateProduct   = "Объявление успешно активировано"
	ResponseSuccessfulAddPremium        = "У объявления успешно акитвирован премиум"
	ResponseSuccessfulDeleteComment     = "Комментарий успешно удалено"
	ResponseSuccessfulUpdateComment     = "Комментарий успешно изменен"
)

//easyjson:json
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

//easyjson:json
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

//easyjson:json
type CommentListResponse struct {
	Status int                     `json:"status"`
	Body   []*models.CommentInFeed `json:"body"`
}

func NewCommentListResponse(body []*models.CommentInFeed) *CommentListResponse {
	return &CommentListResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

//easyjson:json
type ProductInSearchListResponse struct {
	Status int      `json:"status"`
	Body   []string `json:"body"`
}

func NewProductInSearchListResponse(body []string) *ProductInSearchListResponse {
	return &ProductInSearchListResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

//easyjson:json
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

//easyjson:json
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

//easyjson:json
type ResponseAPIYoomany struct {
	Confirmation ConfirmationPayment `json:"confirmation"`
}

func (r *ResponseAPIYoomany) IsCorrect() bool {
	return r.Confirmation.Type == TypeConfirmationPayment
}

//easyjson:json
type OrderNotInBasketListResponse struct {
	Status int                        `json:"status"`
	Body   []*models.OrderNotInBasket `json:"body"`
}

func NewOrderNotInBasketListResponse(body []*models.OrderNotInBasket) *OrderNotInBasketListResponse {
	return &OrderNotInBasketListResponse{Status: statuses.StatusResponseSuccessful, Body: body}
}

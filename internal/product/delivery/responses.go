package delivery

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

const (
	ResponseSuccessfulUpdateCountOrder  = "Успешно изменено количество заказа"
	ResponseSuccessfulUpdateStatusOrder = "Успешно изменен статус заказа"
	ResponseSuccessfulBuyFullBasket     = "Успешная покупка всего из корзины"
	ResponseSuccessfulCloseProduct      = "Объявление успешно закрыто"
	ResponseSuccessfulDeleteProduct     = "Объявление успешно удалено"
	ResponseSuccessfulActivateProduct   = "Объявление успешно активировано"
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
type ConfirmationPayment struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

//easyjson:json
type ResponsePostPaymentAPIYoomany struct {
	Confirmation ConfirmationPayment `json:"confirmation"`
}

func (r *ResponsePostPaymentAPIYoomany) IsCorrect() bool {
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

//easyjson:json
type responseGetPaymentsItemAPIYoomany struct {
	Status   string        `json:"status"`
	Amount   AmountPayment `json:"amount"`
	Metadata struct {
		UserID     string `json:"user_id"`
		ProductID  string `json:"product_id"`
		PeriodCode string `json:"period_code"`
	} `json:"metadata"`
}

//easyjson:json
type responseGetPaymentsAPIYoomany struct {
	Type       string                              `json:"type"`
	NextCursor string                              `json:"next_cursor"`
	Items      []responseGetPaymentsItemAPIYoomany `json:"items"`
}

type ResponseGetPaymentsItemAPIYoomany struct {
	Status   string          `json:"status"`
	Amount   AmountPayment   `json:"amount"`
	Metadata MetadataPayment `json:"metadata"`
}

type ResponseGetPaymentsAPIYoomany struct {
	Type       string                              `json:"type"`
	NextCursor string                              `json:"next_cursor"`
	Items      []ResponseGetPaymentsItemAPIYoomany `json:"items"`
}

func (r *ResponseGetPaymentsAPIYoomany) UnmarshalJSON(body []byte) error {
	var responseGetPaymentsAPIYoomany responseGetPaymentsAPIYoomany

	err := json.Unmarshal(body, &responseGetPaymentsAPIYoomany)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	r.Type = responseGetPaymentsAPIYoomany.Type
	r.NextCursor = responseGetPaymentsAPIYoomany.NextCursor

	for _, item := range responseGetPaymentsAPIYoomany.Items {
		userID, err := strconv.ParseUint(item.Metadata.UserID, 10, 64)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		productID, err := strconv.ParseUint(item.Metadata.ProductID, 10, 64)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		periodCode, err := strconv.ParseUint(item.Metadata.PeriodCode, 10, 64)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		r.Items = append(r.Items, ResponseGetPaymentsItemAPIYoomany{
			Status: item.Status,
			Amount: item.Amount, Metadata: MetadataPayment{
				UserID: userID, ProductID: productID, PeriodCode: periodCode,
			},
		})
	}

	return nil
}

func (r *ResponseGetPaymentsAPIYoomany) MarshalJSON() ([]byte, error) {
	var responseGetPaymentsAPIYoomany responseGetPaymentsAPIYoomany

	responseGetPaymentsAPIYoomany.Type = r.Type
	responseGetPaymentsAPIYoomany.NextCursor = r.NextCursor

	for _, item := range r.Items {
		userID := strconv.FormatUint(item.Metadata.UserID, 10)
		productID := strconv.FormatUint(item.Metadata.ProductID, 10)
		periodCode := strconv.FormatUint(item.Metadata.PeriodCode, 10)

		responseGetPaymentsAPIYoomany.Items = append(responseGetPaymentsAPIYoomany.Items,
			responseGetPaymentsItemAPIYoomany{
				Status: item.Status,
				Amount: item.Amount, Metadata: struct {
					UserID     string `json:"user_id"`
					ProductID  string `json:"product_id"`
					PeriodCode string `json:"period_code"`
				}{
					UserID: userID, ProductID: productID, PeriodCode: periodCode,
				},
			})
	}

	body, err := json.Marshal(responseGetPaymentsAPIYoomany)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return body, nil
}

//easyjson:json
type PremiumStatus struct {
	PremiumStatus uint8 `json:"premium_status"`
}

//easyjson:json
type PremiumStatusResponse struct {
	Status int           `json:"status"`
	Body   PremiumStatus `json:"body"`
}

func NewPremiumStatusResponse(premiumStatus uint8) *PremiumStatusResponse {
	return &PremiumStatusResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   PremiumStatus{PremiumStatus: premiumStatus},
	}
}

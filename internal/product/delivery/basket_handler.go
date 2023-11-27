package delivery

import (
	"context"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ IBasketService = (*productusecases.BasketService)(nil)

type IBasketService interface {
	AddOrder(ctx context.Context, r io.Reader, userID uint64) (*models.OrderInBasket, error)
	GetOrdersByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error)
	UpdateOrderCount(ctx context.Context, r io.Reader, userID uint64) error
	UpdateOrderStatus(ctx context.Context, r io.Reader, userID uint64) error
	BuyFullBasket(ctx context.Context, userID uint64) error
	DeleteOrder(ctx context.Context, orderID uint64, ownerID uint64) error
}

// AddOrderHandler godoc
//
//	@Summary    add order to basket
//	@Description   add product in basket
//	@Tags order
//	@Accept      json
//	@Produce    json
//
// @Param preOrder  body internal_models.PreOrder true  "order data for adding"
//
//	@Success    200  {object} OrderResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /order/add [post]
func (p *ProductHandler) AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	orderInBasket, err := p.service.AddOrder(ctx, r.Body, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendResponse(w, p.logger, NewOrderResponse(orderInBasket))
	p.logger.Infof("in AddOrderHandler: add order orderID=%d for userID=%d\n", orderInBasket.ID, userID)
}

// GetBasketHandler godoc
//
//	@Summary    get basket of orders
//	@Description  get basket of orders by user id from cookie\jwt token
//	@Tags order
//	@Accept     json
//	@Produce    json
//	@Success    200  {object} OrderListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [get]
func (p *ProductHandler) GetBasketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	orders, err := p.service.GetOrdersByUserID(ctx, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendResponse(w, p.logger, NewOrderListResponse(orders))
	p.logger.Infof("in GetBasketHandler: get basket of orders: %+v\n", orders)
}

// UpdateOrderCountHandler godoc
//
//	@Summary    update order count
//	@Description  update order count using user id from cookie\jwt token
//	@Tags order
//	@Accept      json
//	@Produce    json
//
// @Param orderChanges  body internal_models.OrderChanges true  "order data for updating use only id and count"
//
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /order/update_count [patch]
func (p *ProductHandler) UpdateOrderCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.UpdateOrderCount(ctx, r.Body, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendResponse(w, p.logger,
		delivery.NewResponseSuccessful(ResponseSuccessfulUpdateCountOrder))
	p.logger.Infof("in UpdateOrderCountHandler: updated order count for user id=%d\n", userID)
}

// UpdateOrderStatusHandler godoc
//
//	@Summary    update order status
//	@Description  update order status using user id from cookie\jwt token
//	@Tags order
//	@Accept      json
//	@Produce    json
//
// @Param orderChanges  body internal_models.OrderChanges true  "order data for updating use only id and status"
//
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error". Внутри body статус может быть badContent(4400)
//	@Router      /order/update_status [patch]
func (p *ProductHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.UpdateOrderStatus(ctx, r.Body, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendResponse(w, p.logger,
		delivery.NewResponseSuccessful(ResponseSuccessfulUpdateStatusOrder))
	p.logger.Infof("in UpdateOrderStatusHandler: updated order status for user id=%d\n", userID)
}

// BuyFullBasketHandler godoc
//
//	@Summary    buy all orders from basket
//	@Description   buy all orders from basket
//	@Tags order
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error". Внутри body статус может быть badContent(4400)
//	@Router      /order/buy_full_basket [patch]
func (p *ProductHandler) BuyFullBasketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.BuyFullBasket(ctx, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendResponse(w, p.logger,
		delivery.NewResponseSuccessful(ResponseSuccessfulBuyFullBasket))
	p.logger.Infof("in BuyFullBasketHandler: buy full basket for userID=%d\n", userID)
}

// DeleteOrderHandler godoc
//
//	@Summary     delete order
//	@Description  delete order for owner using user id from cookies\jwt.
//	@Description  This totally removed order. Recovery will be impossible
//	@Tags order
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "order id"
//	@Success    200  {object} delivery.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error". Внутри body статус может быть badContent(4400)
//	@Router      /order/delete/ [delete]
func (p *ProductHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	orderID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	err = p.service.DeleteOrder(ctx, orderID, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendResponse(w, p.logger,
		delivery.NewResponseSuccessful(ResponseSuccessfulDeleteProduct))
	p.logger.Infof("in DeleteOrderHandler: delete order id=%d", orderID)
}

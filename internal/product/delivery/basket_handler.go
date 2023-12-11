package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"io"
	"net/http"

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
// @Param preOrder  body models.PreOrder true  "order data for adding"
//
//	@Success    200  {object} OrderResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /order/add [post]
func (p *ProductHandler) AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	orderInBasket, err := p.service.AddOrder(ctx, r.Body, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewOrderResponse(orderInBasket))
	logger.Infof("in AddOrderHandler: add order orderID=%d for userID=%d\n", orderInBasket.ID, userID)
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
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badFormat(4000)
//	@Router      /order/get_basket [get]
func (p *ProductHandler) GetBasketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	orders, err := p.service.GetOrdersByUserID(ctx, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewOrderListResponse(orders))
	logger.Infof("in GetBasketHandler: get basket of orders: %+v\n", orders)
}

// UpdateOrderCountHandler godoc
//
//	@Summary    update order count
//	@Description  update order count using user id from cookie\jwt token
//	@Tags order
//	@Accept      json
//	@Produce    json
//
// @Param orderChanges  body models.OrderChanges true  "order data for updating use only id and count"
//
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400), badFormat(4000)
//	@Router      /order/update_count [patch]
func (p *ProductHandler) UpdateOrderCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.UpdateOrderCount(ctx, r.Body, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulUpdateCountOrder))
	logger.Infof("in UpdateOrderCountHandler: updated order count for user id=%d\n", userID)
}

// UpdateOrderStatusHandler godoc
//
//	@Summary    update order status
//	@Description  update order status using user id from cookie\jwt token
//	@Tags order
//	@Accept      json
//	@Produce    json
//
// @Param orderChanges  body models.OrderChanges true  "order data for updating use only id and status"
//
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400)
//	@Router      /order/update_status [patch]
func (p *ProductHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.UpdateOrderStatus(ctx, r.Body, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulUpdateStatusOrder))
	logger.Infof("in UpdateOrderStatusHandler: updated order status for user id=%d\n", userID)
}

// BuyFullBasketHandler godoc
//
//	@Summary    buy all orders from basket
//	@Description   buy all orders from basket
//	@Tags order
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400)
//	@Router      /order/buy_full_basket [patch]
func (p *ProductHandler) BuyFullBasketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.BuyFullBasket(ctx, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulBuyFullBasket))
	logger.Infof("in BuyFullBasketHandler: buy full basket for userID=%d\n", userID)
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
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badContent(4400)
//	@Router      /order/delete/ [delete]
func (p *ProductHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	orderID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.DeleteOrder(ctx, orderID, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulDeleteProduct))
	logger.Infof("in DeleteOrderHandler: delete order id=%d", orderID)
}

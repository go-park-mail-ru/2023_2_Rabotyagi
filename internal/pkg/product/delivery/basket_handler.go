package delivery

import (
	"net/http"
	"strconv"

	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
)

// GetBasketHandler godoc
//
//	@Summary    get basket of orders
//	@Description  get basket of orders by user id from cookie\jwt token
//	@Tags order
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} OrderListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [get]
func (p *ProductHandler) GetBasketHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r, p.logger)

	orders, err := p.storage.GetOrdersInBasketByUserID(ctx, userID)
	if err != nil {
		p.logger.Errorf("in GetBasketHandler %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	for _, order := range orders {
		order.Sanitize()
	}

	delivery.SendOkResponse(w, p.logger, NewOrderListResponse(delivery.StatusResponseSuccessful, orders))
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
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/update_count [patch]
func (p *ProductHandler) UpdateOrderCountHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	orderChanges, err := productusecases.ValidateOrderChangesCount(p.logger, r.Body)
	if err != nil {
		p.logger.Errorf("in in UpdateOrderCountHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	userID := delivery.GetUserIDFromCookie(r, p.logger)

	err = p.storage.UpdateOrderCount(ctx, userID, orderChanges.ID, orderChanges.Count)
	if err != nil {
		p.logger.Errorf("in in UpdateOrderCountHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulUpdateCountOrder))
	p.logger.Infof("in UpdateOrderCountHandler: updated order count=%d for order id=%d for user id=%d\n",
		orderChanges.Count, orderChanges.ID, userID)
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
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/update_status [patch]
func (p *ProductHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	orderChanges, err := productusecases.ValidateOrderChangesStatus(p.logger, r.Body)
	if err != nil {
		p.logger.Errorf("in UpdateOrderStatusHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r, p.logger)

	err = p.storage.UpdateOrderStatus(ctx, userID, orderChanges.ID, orderChanges.Status)
	if err != nil {
		p.logger.Errorf("in UpdateOrderStatusHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulUpdateStatusOrder))
	p.logger.Infof("in UpdateOrderStatusHandler: updated order id=%d with status=%d for user id=%d\n",
		orderChanges.ID, orderChanges.Status, userID)
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/add [post]
func (p *ProductHandler) AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r, p.logger)

	preOrder, err := productusecases.ValidatePreOrder(p.logger, r.Body)
	if err != nil {
		p.logger.Errorf("in AddOrderHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	orderInBasket, err := p.storage.AddOrderInBasket(ctx, userID, preOrder.ProductID, preOrder.Count)
	if err != nil {
		p.logger.Errorf("in AddOrderHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewOrderResponse(delivery.StatusResponseSuccessful, orderInBasket))
	p.logger.Infof("in AddOrderHandler: add order on productID=%d for userID=%d\n", preOrder.ProductID, userID)
}

// BuyFullBasketHandler godoc
//
//	@Summary    buy all orders from basket
//	@Description   buy all orders from basket
//	@Tags order
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/buy_full_basket [patch]
func (p *ProductHandler) BuyFullBasketHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r, p.logger)

	err := p.storage.BuyFullBasket(ctx, userID)
	if err != nil {
		p.logger.Errorf("in BuyFullBasketHandler: %+v\n", err)
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulBuyFullBasket))
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
//	@Param      orderID  path uint64 true  "order id"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/delete/ [delete]
func (p *ProductHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r, p.logger)
	orderIDStr := delivery.GetPathParam(r.URL.String())

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		p.logger.Errorf("in DeleteOrderHandler: %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrWrongProductID.Error()))

		return
	}

	err = p.storage.DeleteOrder(ctx, orderID, userID)
	if err != nil {
		p.logger.Errorf("in DeleteOrderHandler %+v\n", err)
		delivery.SendErrResponse(w, p.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, p.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulDeleteProduct))
	p.logger.Infof("in DeleteOrderHandler: delete order id=%d", orderID)
}

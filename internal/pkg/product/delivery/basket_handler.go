package delivery

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
)

// GetBasketHandler godoc
//
//	@Summary    get basket of orders
//	@Description  get basket of orders by user id from cookie\jwt token
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} OrderListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [get]
func (p *ProductHandler) GetBasketHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r)

	orders, err := p.storage.GetOrdersInBasketByUserID(ctx, userID)
	if err != nil {
		log.Printf("in GetBasketHandler %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delivery.SendOkResponse(w, NewOrderListResponse(delivery.StatusResponseSuccessful, orders))
	log.Printf("in GetBasketHandler: get basket of orders: %+v\n", orders)
}

// UpdateOrderCountHandler godoc
//
//	@Summary    update order count
//	@Description  update order count using user id from cookie\jwt token
//	@Accept      json
//	@Produce    json
//
// @Param orderChanges  body internal_models.OrderChanges true  "order data for updating use only id and count"
//
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [patch]
func (p *ProductHandler) UpdateOrderCountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	orderChanges, err := usecases.ValidateOrderChangesCount(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderCountHandler:", err)

		return
	}

	userID := delivery.GetUserIDFromCookie(r)

	err = p.storage.UpdateOrderCount(ctx, userID, orderChanges.ID, orderChanges.Count)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderCountHandler:", err)

		return
	}

	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulUpdateCountOrder))
	log.Printf("in UpdateOrderCountHandler: updated order count=%d for order id=%d for user id=%d\n",
		orderChanges.Count, orderChanges.ID, userID)
}

// UpdateOrderStatusHandler godoc
//
//	@Summary    update order status
//	@Description  update order status using user id from cookie\jwt token
//	@Accept      json
//	@Produce    json
//
// @Param orderChanges  body internal_models.OrderChanges true  "order data for updating use only id and status"
//
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [patch]
func (p *ProductHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	orderChanges, err := usecases.ValidateOrderChangesStatus(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderStatusHandler:", err)

		return
	}

	ctx := r.Context()
	userID := delivery.GetUserIDFromCookie(r)

	err = p.storage.UpdateOrderStatus(ctx, userID, orderChanges.ID, orderChanges.Status)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderStatusHandler:", err)

		return
	}

	delivery.SendOkResponse(w,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulUpdateStatusOrder))
	log.Printf("in UpdateOrderStatusHandler: updated order id=%d with status=%d for user id=%d\n",
		orderChanges.ID, orderChanges.Status, userID)
}

// AddOrderHandler godoc
//
//	@Summary    add order to basket
//	@Description   add product in basket
//	@Accept      json
//	@Produce    json
//
// @Param preOrder  body internal_models.PreOrder true  "order data for adding"
//
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/add [post]
func (p *ProductHandler) AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r)

	preOrder, err := usecases.ValidatePreOrder(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in AddOrderHandler:", err)

		return
	}

	err = p.storage.AddOrderInBasket(ctx, userID, preOrder.ProductID, preOrder.Count)
	if err != nil {
		delivery.HandleErr(w, "in AddOrderHandler:", err)

		return
	}

	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulAddOrder))
	log.Printf("in AddOrderHandler: add order on productID=%d for userID=%d\n", preOrder.ProductID, userID)
}

// BuyFullBasketHandler godoc
//
//	@Summary    buy all orders from basket
//	@Description   buy all orders from basket
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/add [patch]
func (p *ProductHandler) BuyFullBasketHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID := delivery.GetUserIDFromCookie(r)

	err := p.storage.BuyFullBasket(ctx, userID)
	if err != nil {
		delivery.HandleErr(w, "in BuyFullBasketHandler:", err)

		return
	}

	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulBuyFullBasket))
	log.Printf("in BuyFullBasketHandler: buy full basket for userID=%d\n", userID)
}

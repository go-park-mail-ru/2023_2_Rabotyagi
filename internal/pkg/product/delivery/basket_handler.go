package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"

	"github.com/asaskevich/govalidator"
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
	log.Printf("in GetBasketHandler: get basket of orders: %+v", orders)
}

// UpdateOrderCountHandler godoc
//
//	@Summary    update order count
//	@Description  update order count using user id from cookie\jwt token
//	@Accept      json
//	@Produce    json
//
// @Param count  body internal_models.OrderChanges true  "order data for updating"
//
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [path]
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
	log.Printf("in UpdateOrderCountHandler: updated order count=%d for order id=%d for user id=%d:",
		orderChanges.Count, orderChanges.ID, userID)
}

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

	ctx := r.Context()

	newOrder := struct {
		ID     uint64 `json:"id"          valid:"required"`
		Status uint8  `json:"status"       valid:"required"`
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newOrder); err != nil {
		delivery.HandleErr(w, "in UpdateOrderStatusHandler:", err)

		return
	}

	_, err := govalidator.ValidateStruct(newOrder)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderStatusHandler:", err)

		return
	}

	updatedOrder, err := p.storage.UpdateOrderStatus(ctx, newOrder.ID, newOrder.Status)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderStatusHandler:", err)

		return
	}

	delivery.SendOkResponse(w, NewOrderResponse(delivery.StatusResponseSuccessful, updatedOrder))
	log.Printf("in UpdateOrderStatusHandler: update order status: %+v", updatedOrder)
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
	log.Printf("in AddOrderHandler: add order on productID=%d for userID=%d", preOrder.ProductID, userID)
}

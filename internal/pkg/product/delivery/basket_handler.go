package delivery

import (
	"encoding/json"
	"log"
	"net/http"

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
//	@Summary    get basket of orders
//	@Description  get basket of orders by user id from cookie\jwt token
//	@Accept      json
//	@Produce    json
//	@Success    200  {object} OrderListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /order/get_basket [get]
func (p *ProductHandler) UpdateOrderCountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	newOrder := struct {
		ID    uint64 `json:"id"          valid:"required"`
		Count uint32 `json:"count"       valid:"required"`
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newOrder); err != nil {
		delivery.HandleErr(w, "in UpdateOrderCountHandler:", err)

		return
	}

	_, err := govalidator.ValidateStruct(newOrder)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderCountHandler:", err)

		return
	}

	updatedOrder, err := p.storage.UpdateOrderCount(ctx, newOrder.ID, newOrder.Count)
	if err != nil {
		delivery.HandleErr(w, "in UpdateOrderCountHandler:", err)

		return
	}

	delivery.SendOkResponse(w, NewOrderResponse(delivery.StatusResponseSuccessful, updatedOrder))
	log.Printf("in UpdateOrderCountHandler: update order count: %+v", updatedOrder)
}

func (p *ProductHandler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.addrOrigin, p.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
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

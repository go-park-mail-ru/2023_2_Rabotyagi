package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"net/http"
)

var _ IFavouriteService = (*productusecases.FavouriteService)(nil)

type IFavouriteService interface {
	GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error)
	AddToFavourites(ctx context.Context, userID uint64, productID uint64) error
	DeleteFromFavourites(ctx context.Context, userID uint64, productID uint64) error
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
func (p *ProductHandler) GetFavouritesHandler(w http.ResponseWriter, r *http.Request) {
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

	products, err := p.service.Ge(ctx, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewOrderListResponse(delivery.StatusResponseSuccessful, orders))
	p.logger.Infof("in GetBasketHandler: get basket of orders: %+v\n", orders)
}

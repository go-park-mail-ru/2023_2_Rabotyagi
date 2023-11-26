package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"net/http"
	"strconv"
)

var _ IFavouriteService = (*productusecases.FavouriteService)(nil)

type IFavouriteService interface {
	GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error)
	AddToFavourites(ctx context.Context, userID uint64, productID uint64) error
	DeleteFromFavourites(ctx context.Context, userID uint64, productID uint64) error
}

// GetFavouritesHandler godoc
//
//	@Summary    get user favourites
//	@Description  get user favourites by user id from cookie\jwt token
//	@Tags favourite
//	@Accept     json
//	@Produce    json
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/favourites [get]
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

	products, err := p.service.GetUserFavourites(ctx, userID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, NewProductListResponse(delivery.StatusResponseSuccessful, products))
	p.logger.Infof("in GetFavouritesHandler: get user favourites: %+v\n", products)
}

// AddToFavouritesHandler godoc
//
//	@Summary     add product to favs
//	@Description  add product to favs using product id from query and user id form cookie
//	@Tags favourite
//	@Accept      json
//	@Produce    json
//	@Param      product_id  query uint64 true  "product id"
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/add-to-fav [post]
func (p *ProductHandler) AddToFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	productIDStr := r.URL.Query().Get("product_id")

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	err = p.service.AddToFavourites(ctx, userID, productID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, delivery.NewResponseID(productID))
	p.logger.Infof("in AddToFavouritesHandler: add to fav with product id = %+v", productID)
}

// DeleteFromFavouritesHandler godoc
//
//	@Summary     delete product from favs
//	@Description  delete product from favs using product id from query and user id form cookie
//	@Tags favourite
//	@Accept      json
//	@Produce    json
//	@Param      product_id  query uint64 true  "product id"
//	@Success    200  {object} ProductListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/remove-from-fav [delete]
func (p *ProductHandler) DeleteFromFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	productIDStr := r.URL.Query().Get("product_id")

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	ctx := r.Context()

	err = p.service.DeleteFromFavourites(ctx, userID, productID)
	if err != nil {
		delivery.HandleErr(w, p.logger, err)

		return
	}

	delivery.SendOkResponse(w, p.logger, delivery.NewResponseID(productID))
	p.logger.Infof("in DeleteFromFavouritesHandler: del form fav with product id = %+v", productID)
}

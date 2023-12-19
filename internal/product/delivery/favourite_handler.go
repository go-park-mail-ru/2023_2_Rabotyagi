package delivery

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"io"
	"net/http"
	"strconv"

	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
)

var _ IFavouriteService = (*productusecases.FavouriteService)(nil)

type IFavouriteService interface {
	GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error)
	AddToFavourites(ctx context.Context, userID uint64, r io.Reader) error
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
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badFormat(4000)
//	@Router      /profile/favourites [get]
func (p *ProductHandler) GetFavouritesHandler(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
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

	products, err := p.service.GetUserFavourites(ctx, userID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger, NewProductListResponse(products))
	logger.Infof("in GetFavouritesHandler: get user favourites: %+v\n", products)
}

// AddToFavouritesHandler godoc
//
//	@Summary     add product to favs
//	@Description  add product to favs using product id from query and user id form cookie
//	@Tags favourite
//	@Accept      json
//	@Produce    json
//	@Param      product_id  body models.ProductID true "product id"
//	@Success    200  {object} responses.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badFormat(4000)
//	@Router      /product/add-to-fav [post]
func (p *ProductHandler) AddToFavouritesHandler(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
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

	err = p.service.AddToFavourites(ctx, userID, r.Body)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	r = r.WithContext(statuses.FillStatusCtx(r.Context(), statuses.StatusRedirectAfterSuccessful))

	responses.SendResponse(w, logger, responses.NewResponseIDRedirect(userID))

	logger.Infof("in AddToFavouritesHandler: add to fav with user id = %+v", userID)
}

// DeleteFromFavouritesHandler godoc
//
//	@Summary     delete product from favs
//	@Description  delete product from favs using product id from query and user id form cookie
//	@Tags favourite
//	@Accept      json
//	@Produce    json
//	@Param      product_id  query uint64 true  "product id"
//	@Success    200  {object} responses.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Внутри body статус может быть badFormat(4000)
//
// @Router /product/remove-from-fav [delete]
func (p *ProductHandler) DeleteFromFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	productIDStr := r.URL.Query().Get("product_id")

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.DeleteFromFavourites(ctx, userID, productID)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	r = r.WithContext(statuses.FillStatusCtx(r.Context(), statuses.StatusRedirectAfterSuccessful))

	responses.SendResponse(w, logger, responses.NewResponseIDRedirect(productID))

	logger.Infof("in DeleteFromFavouritesHandler: del form fav with product id = %+v", productID)
}

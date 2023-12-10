package delivery

import (
	"context"
	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"net/http"
)

const (
	ResponseSuccessfulAddPremium = "У объявления успешно акитвирован премиум"
)

var _ IPremiumService = (*productusecases.PremiumService)(nil)

type IPremiumService interface {
	AddPremium(ctx context.Context, productID uint64, userID uint64) error
	RemovePremium(ctx context.Context, productID uint64, userID uint64) error
}

// AddPremiumHandler godoc
//
//	@Summary     add premium for product
//	@Description  add premium for product using product id from query and user id from cookies\jwt.
//	@Description  This does product active.
//	@Tags premium
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "product id"
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badFormat(4000)
//	@Router      /product/premium/add [patch]
func (p *ProductHandler) AddPremiumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
	logger := p.logger.LogReqID(ctx)

	userID, err := delivery.GetUserID(ctx, r, p.sessionManagerClient)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	productID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	err = p.service.AddPremium(ctx, productID, userID)
	if err != nil {
		responses.HandleErr(w, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseSuccessful(ResponseSuccessfulAddPremium))
	logger.Infof("in AddPremiumHandler: product id=%d", productID)
}

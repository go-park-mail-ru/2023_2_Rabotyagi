package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	productusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var _ IPremiumService = (*productusecases.PremiumService)(nil)

type IPremiumService interface {
	AddPremium(ctx context.Context, productID uint64, userID uint64, periodCode uint64) error
}

var (
	ErrMarshallPayment           = myerrors.NewErrorInternal("Ошибка маршалинга платежа")
	ErrCreationRequestAPIYooMany = myerrors.NewErrorInternal("Ошибка создания запроса к yoomany")
	ErrRequestAPIYoomany         = myerrors.NewErrorInternal("Ошибка в заросе к yoomany")
	ErrReadAllAPIYoomany         = myerrors.NewErrorInternal("Ошибка в чтении ответа от yoomany")
	ErrUnmarshallAPIYoomany      = myerrors.NewErrorInternal("Ошибка в unmarshall от yoomany")
	ErrResponseAPIYoomany        = myerrors.NewErrorInternal("Ошибка проверки ответа от yoomany")
)

const (
	headerKeyIdempotency  = "Idempotency-Key"
	yoomanyPaymentsAPIURL = "https://api.yookassa.ru/v3/payments"
)

//nolint:funlen
func (p *ProductHandler) createPayment(ctx context.Context,
	userID uint64, productID uint64, periodCode uint64,
) (string, error) {
	logger := p.logger.LogReqID(ctx)

	payment, err := NewPayment(ctx, p.frontendURL, NewMetadataPayment(userID, productID, periodCode))
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	body, err := payment.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrMarshallPayment, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	keyIdempotencyPayment := p.mapIdempotencyPayment.AddPayment(payment.Metadata)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, yoomanyPaymentsAPIURL, bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrCreationRequestAPIYooMany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	req.Header.Set(headerKeyIdempotency, string(keyIdempotencyPayment))
	req.SetBasicAuth(p.premiumShopID, p.premiumShopSecretKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrRequestAPIYoomany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	defer response.Body.Close()

	bodyResp, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrReadAllAPIYoomany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var responsePayment ResponseAPIYoomany

	err = json.Unmarshal(bodyResp, &responsePayment)
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrUnmarshallAPIYoomany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if !responsePayment.IsCorrect() {
		logger.Errorln(ErrResponseAPIYoomany)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return responsePayment.Confirmation.ReturnURL, nil
}

// AddPremiumHandler godoc
//
//	@Summary     add premium for product
//	@Description  add premium for product using product id from query and user id from cookies\jwt.
//	@Description  This does product premium.
//	@Tags premium
//	@Accept      json
//	@Produce    json
//	@Param      product_id  query uint64 true  "product id"
//	@Param      period  query uint64 true  "period of premium"
//	@Success    200  {object} responses.ResponseRedirect
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error". Это Http ответ 200, внутри body статус может быть badFormat(4000)//nolint:lll
//	@Router      /premium/add [patch]
func (p *ProductHandler) AddPremiumHandler(w http.ResponseWriter, r *http.Request) {
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

	productID, err := utils.ParseUint64FromRequest(r, "product_id")
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	periodCode, err := utils.ParseUint64FromRequest(r, "period")
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	redirectURL, err := p.createPayment(ctx, userID, productID, periodCode)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	err = p.service.AddPremium(ctx, productID, userID, periodCode)
	if err != nil {
		responses.HandleErr(w, r, logger, err)

		return
	}

	responses.SendResponse(w, logger,
		responses.NewResponseRedirect(redirectURL))
	logger.Infof("in AddPremiumHandler: product id=%d userID=%d periodCode=%d", productID, userID, periodCode)
}

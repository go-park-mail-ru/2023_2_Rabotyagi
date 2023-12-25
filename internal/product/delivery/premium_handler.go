package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

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
	ErrMarshallPayment               = myerrors.NewErrorInternal("Ошибка маршалинга платежа")
	ErrCreationRequestAPIYooMany     = myerrors.NewErrorInternal("Ошибка создания запроса к yoomany")
	ErrClosingResponseBody           = myerrors.NewErrorInternal("Ошибка закрытия тела ответа")
	ErrRequestAPIYoomany             = myerrors.NewErrorInternal("Ошибка в заросе к yoomany")
	ErrReadAllAPIYoomany             = myerrors.NewErrorInternal("Ошибка в чтении ответа от yoomany")
	ErrUnmarshallAPIYoomany          = myerrors.NewErrorInternal("Ошибка в unmarshall от yoomany")
	ErrResponseAPIYoomany            = myerrors.NewErrorInternal("Ошибка проверки ответа от yoomany")
	ErrResponseWrongStatusAPIYoomany = myerrors.NewErrorBadContentRequest("Ошибка оплаты платежа от yoomany")
	ErrDidntWaitPaymentAPIYoomany    = myerrors.NewErrorBadContentRequest("Не дождались оплаты")
	ErrValidationPaymentAPIYoomany   = myerrors.NewErrorInternal("Оплата не прошла валидацию")
	ErrNotFoundPaymentAPIYoomany     = myerrors.NewErrorInternal("Не найден нужный платеж")
	ErrPaymentPaindingAPIYoomany     = myerrors.NewErrorInternal("Платеж в обработке")
)

const (
	headerKeyIdempotency     = "Idempotence-Key"
	paymentsURLAPIYoomany    = "https://api.yookassa.ru/v3/payments"
	paramCreatedAtAPIYoomany = "created_at.gte="
	maxTimeoutAPIYoumany     = time.Minute * 11
	periodRequestAPIYoumany  = time.Second * 11
)

// parsePayments true in return argument means successful handle payment
//
//nolint:funlen,cyclop
func (p *ProductHandler) parsePayments(ctx context.Context, payment *Payment, reader io.Reader) (bool, error) {
	logger := p.logger.LogReqID(ctx)

	body, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("%w %+v", ErrReadAllAPIYoomany, err) //nolint:errorlint
		logger.Errorln(err)

		return false, err
	}

	logger.Infof("body:%s", body)

	var responseGetPayments ResponseGetPaymentsAPIYoomany

	err = json.Unmarshal(body, &responseGetPayments)
	if err != nil {
		err = fmt.Errorf("%w %+v", ErrUnmarshallAPIYoomany, err) //nolint:errorlint
		logger.Errorln(err)

		return false, err
	}

	for _, item := range responseGetPayments.Items {
		logger.Infof("%+v\n%+v", item, payment)

		if item.Metadata.UserID == payment.Metadata.UserID &&
			item.Metadata.ProductID == payment.Metadata.ProductID &&
			item.Metadata.PeriodCode == payment.Metadata.PeriodCode {
			switch {
			case item.Status == StatusPaymentPending:
				logger.Errorln(ErrPaymentPaindingAPIYoomany)

				return false, nil
			case IsStatusPaymentSuccessful(item.Status):
				if !reflect.DeepEqual(item.Amount, payment.Amount) {
					err = fmt.Errorf("%w received: %+v != requested: %+v",
						ErrValidationPaymentAPIYoomany, item.Amount, payment.Amount)

					logger.Errorln(err)

					return false, err
				}

				err = p.service.AddPremium(ctx,
					payment.Metadata.ProductID, payment.Metadata.UserID, payment.Metadata.PeriodCode)
				if err != nil {
					err = fmt.Errorf(myerrors.ErrTemplate, err)
					logger.Errorln(err)

					return false, err
				}

				logger.Infof("Successful addPremium metadata:%+v", payment.Metadata)

				return true, nil
			default:
				logger.Errorln(ErrResponseWrongStatusAPIYoomany)

				return false, fmt.Errorf(myerrors.ErrTemplate, ErrResponseWrongStatusAPIYoomany)
			}
		}
	}

	logger.Errorln(ErrNotFoundPaymentAPIYoomany)

	return false, nil
}

func (p *ProductHandler) waitPayment(ctx context.Context, chError chan<- error,
	payment *Payment, periodRequest time.Duration,
) {
	logger := p.logger.LogReqID(ctx)
	timer := time.NewTimer(maxTimeoutAPIYoumany)
	timeStart := time.Now().Add(-100 * time.Second).Format(time.RFC3339)

	go func() {
		for {
			select {
			case <-timer.C:
				err := fmt.Errorf("%w для %+v", ErrDidntWaitPaymentAPIYoomany, payment)

				logger.Errorln(err)
				chError <- err
			default:
				time.Sleep(periodRequest)

				request, err := http.NewRequestWithContext(ctx,
					http.MethodGet, fmt.Sprintf("%s?%s%s", paymentsURLAPIYoomany, paramCreatedAtAPIYoomany, timeStart), nil)
				if err != nil {
					err = fmt.Errorf("%w %+v", ErrCreationRequestAPIYooMany, err) //nolint:errorlint
					logger.Errorln(err)
					chError <- err
				}

				request.SetBasicAuth(p.premiumShopID, p.premiumShopSecretKey)
				logger.Infof("req:%+v", request)

				response, err := p.httpClient.Do(request)
				if err != nil {
					err = fmt.Errorf("%w %+v", ErrRequestAPIYoomany, err) //nolint:errorlint
					logger.Errorln(err)
					chError <- err
				}

				isSuccessfully, errPayments := p.parsePayments(ctx, payment, response.Body)

				errBodyClose := response.Body.Close()
				if errBodyClose != nil {
					err = fmt.Errorf("%w %+v", ErrClosingResponseBody, err) //nolint:errorlint
					logger.Errorln(err)
					chError <- err
				}

				if errPayments != nil {
					chError <- err

					return
				}

				if errPayments == nil && isSuccessfully {
					return
				}
			}
		}
	}()
}

//nolint:funlen
func (p *ProductHandler) createPayment(ctx context.Context,
	userID uint64, productID uint64, periodCode uint64,
) (string, error) {
	logger := p.logger.LogReqID(ctx)

	payment, err := NewPayment(ctx, p.frontendPaymentURL, NewMetadataPayment(userID, productID, periodCode))
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	body, err := payment.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrMarshallPayment, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger.Infof("%s", body)

	keyIdempotencyPayment := p.mapIdempotencyPayment.AddPayment(payment.Metadata)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, paymentsURLAPIYoomany, bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrCreationRequestAPIYooMany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	req.Header.Set(headerKeyIdempotency, string(keyIdempotencyPayment))
	req.SetBasicAuth(p.premiumShopID, p.premiumShopSecretKey)
	req.Header.Set("Content-Type", "application/json")
	logger.Infof("%+v", req)

	response, err := p.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrRequestAPIYoomany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger.Infof("%+v", response)

	defer response.Body.Close()

	bodyResp, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("%w error:%v", ErrReadAllAPIYoomany, err.Error())
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var responsePayment ResponsePostPaymentAPIYoomany

	logger.Infof("%s", bodyResp)

	err = json.Unmarshal(bodyResp, &responsePayment)
	if err != nil {
		err = fmt.Errorf("%w error:%+v response: %s", ErrUnmarshallAPIYoomany, err.Error(), bodyResp)
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if !responsePayment.IsCorrect() {
		logger.Errorln(ErrResponseAPIYoomany)
		logger.Infof("response Confirmation %+v", responsePayment.Confirmation)
		logger.Infof("expected Confirmation %+v", payment.Confirmation)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	//nolint:godox
	// TODO chErr just don`t handle yet
	chErr := make(chan error)
	p.waitPayment(ctx, chErr,
		payment, periodRequestAPIYoumany)

	return responsePayment.Confirmation.ConfirmationURL, nil
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

	responses.SendResponse(w, logger,
		responses.NewResponseRedirect(redirectURL))
	logger.Infof("in AddPremiumHandler: product id=%d userID=%d periodCode=%d", productID, userID, periodCode)
}

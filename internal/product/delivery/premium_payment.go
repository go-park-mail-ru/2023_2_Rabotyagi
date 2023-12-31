package delivery

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
)

const pathRedirectURLPremium = "/profile/products"

//easyjson:json
type MetadataPayment struct {
	UserID     uint64 `json:"user_id"`
	ProductID  uint64 `json:"product_id"`
	PeriodCode uint64 `json:"period_code"`
}

func NewMetadataPayment(userID uint64, productID uint64, periodCode uint64) *MetadataPayment {
	return &MetadataPayment{UserID: userID, ProductID: productID, PeriodCode: periodCode}
}

const currencyAmountPayment = "RUB"

//easyjson:json
type AmountPayment struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

func NewAmountPayment(value string) AmountPayment {
	return AmountPayment{Currency: currencyAmountPayment, Value: value}
}

const TypeConfirmationPayment = "redirect"

//easyjson:json
type ConfirmationReturnPayment struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

func NewConfirmationReturnPayment(returnURL string) ConfirmationReturnPayment {
	return ConfirmationReturnPayment{Type: TypeConfirmationPayment, ReturnURL: returnURL}
}

//easyjson:json
type Payment struct {
	Amount       AmountPayment             `json:"amount"`
	Capture      bool                      `json:"capture"`
	Confirmation ConfirmationReturnPayment `json:"confirmation"`
	Description  string                    `json:"description"`
	Metadata     *MetadataPayment          `json:"metadata"`
}

func NewPayment(ctx context.Context, frontendURL string, metadata *MetadataPayment) (*Payment, error) {
	description, err := usecases.GenerateDescriptionByPeriodCode(ctx, metadata.PeriodCode)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	amount, err := usecases.GenerateAmountByPeriodCode(ctx, metadata.PeriodCode)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &Payment{
		Amount:       NewAmountPayment(amount),
		Capture:      true,
		Confirmation: NewConfirmationReturnPayment("https://" + frontendURL + pathRedirectURLPremium),
		Description:  description,
		Metadata:     metadata,
	}, nil
}

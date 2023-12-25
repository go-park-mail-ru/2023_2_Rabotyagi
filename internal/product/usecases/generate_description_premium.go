package usecases

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

const (
	DescriptionWeek       = "Платное продвижение объявления на неделю"
	DescriptionMonth      = "Платное продвижение объявления на 1 месяц"
	DescriptionThreeMonth = "Платное продвижение объявления на 3 месяца"
	DescriptionHalfYear   = "Платное продвижение объявления на 6 месяцев"
	DescriptionYear       = "Платное продвижение объявления на 1 год"

	AmountWeek       = "100.00"
	AmountMonth      = "300.00"
	AmountThreeMonth = "800.00"
	AmountHalfYear   = "1200.00"
	AmountYear       = "1500.00"
)

func GenerateDescriptionByPeriodCode(ctx context.Context, periodCode uint64) (string, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger = logger.LogReqID(ctx)

	switch periodCode {
	case Week:
		return DescriptionWeek, nil
	case Month:
		return DescriptionMonth, nil
	case ThreeMonth:
		return DescriptionThreeMonth, nil
	case HalfYear:
		return DescriptionHalfYear, nil
	case Year:
		return DescriptionYear, nil
	default:
		logger.Errorln(ErrPremiumCode)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrPremiumCode)
	}
}

func GenerateAmountByPeriodCode(ctx context.Context, periodCode uint64) (string, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger = logger.LogReqID(ctx)

	switch periodCode {
	case Week:
		return AmountWeek, nil
	case Month:
		return AmountMonth, nil
	case ThreeMonth:
		return AmountThreeMonth, nil
	case HalfYear:
		return AmountHalfYear, nil
	case Year:
		return AmountYear, nil
	default:
		logger.Errorln(ErrPremiumCode)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrPremiumCode)
	}
}

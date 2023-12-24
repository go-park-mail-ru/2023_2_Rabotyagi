package usecases

import (
	"context"
	"fmt"
	"time"

	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

var ErrPremiumCode = myerrors.NewErrorBadFormatRequest("Ошибка срока подписки на премиум ")

const (
	Week       = uint64(1)
	Month      = uint64(2)
	ThreeMonth = uint64(3)
	HalfYear   = uint64(4)
	Year       = uint64(5)

	DaysInWeek      = 7
	MonthInSeason   = 3
	MonthInHalfYear = 6
)

var _ IPremiumStorage = (*productrepo.ProductStorage)(nil)

type IPremiumStorage interface {
	AddPremium(ctx context.Context, productID uint64, userID uint64,
		premiumBegin time.Time, premiumExpire time.Time) error
}

type PremiumService struct {
	storage IPremiumStorage
	logger  *mylogger.MyLogger
}

func NewPremiumService(premiumStorage IPremiumStorage) (*PremiumService, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &PremiumService{storage: premiumStorage, logger: logger}, nil
}

func (p PremiumService) AddPremium(ctx context.Context, productID uint64,
	userID uint64, periodCode uint64,
) error {
	var premiumExpire time.Time

	premiumBegin := time.Now()
	logger := p.logger.LogReqID(ctx)

	switch periodCode {
	case Week:
		premiumExpire = premiumBegin.AddDate(0, 0, DaysInWeek)
	case Month:
		premiumExpire = premiumBegin.AddDate(0, 1, 0)
	case ThreeMonth:
		premiumExpire = premiumBegin.AddDate(0, MonthInSeason, 0)
	case HalfYear:
		premiumExpire = premiumBegin.AddDate(0, MonthInHalfYear, 0)
	case Year:
		premiumExpire = premiumBegin.AddDate(1, 0, 0)
	default:
		logger.Error(ErrPremiumCode)

		return fmt.Errorf(myerrors.ErrTemplate, ErrPremiumCode)
	}

	err := p.storage.AddPremium(ctx, productID, userID, premiumBegin, premiumExpire)
	if err != nil {
		logger.Error(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

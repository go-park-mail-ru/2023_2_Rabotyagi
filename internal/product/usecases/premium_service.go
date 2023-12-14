package usecases

import (
	"context"
	"fmt"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"time"
)

var (
	ErrPremiumCode = myerrors.NewErrorBadFormatRequest("Ошибка срока подписки на премиум ")
)

const (
	Week       = uint64(1)
	Month      = uint64(2)
	ThreeMonth = uint64(3)
	HalfYear   = uint64(4)
	Year       = uint64(5)
)

var _ IPremiumStorage = (*productrepo.ProductStorage)(nil)

type IPremiumStorage interface {
	AddPremium(ctx context.Context, productID uint64, userID uint64,
		premiumBegin time.Time, premiumExpire time.Time) error
	RemovePremium(ctx context.Context, productID uint64, userID uint64) error
}

type PremiumService struct {
	storage IPremiumStorage
	logger  *my_logger.MyLogger
}

func NewPremiumService(premiumStorage IPremiumStorage) (*PremiumService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &PremiumService{storage: premiumStorage, logger: logger}, nil
}

func (p PremiumService) AddPremium(ctx context.Context, productID uint64,
	userID uint64, periodCode uint64) error {
	var premiumExpire time.Time

	premiumBegin := time.Now()

	switch periodCode {
	case Week:
		premiumExpire = premiumBegin.AddDate(0, 0, 7)
	case Month:
		premiumExpire = premiumBegin.AddDate(0, 1, 0)
	case ThreeMonth:
		premiumExpire = premiumBegin.AddDate(0, 3, 0)
	case HalfYear:
		premiumExpire = premiumBegin.AddDate(0, 6, 0)
	case Year:
		premiumExpire = premiumBegin.AddDate(1, 0, 0)
	default:
		return fmt.Errorf(myerrors.ErrTemplate, ErrPremiumCode)
	}

	err := p.storage.AddPremium(ctx, productID, userID, premiumBegin, premiumExpire)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p PremiumService) RemovePremium(ctx context.Context, productID uint64, userID uint64) error {
	err := p.storage.RemovePremium(ctx, productID, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

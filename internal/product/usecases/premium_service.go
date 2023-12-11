package usecases

import (
	"context"
	"fmt"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
)

var _ IPremiumStorage = (*productrepo.ProductStorage)(nil)

type IPremiumStorage interface {
	AddPremium(ctx context.Context, productID uint64, userID uint64) error
	RemovePremium(ctx context.Context, productID uint64, userID uint64) error
}

type PremiumService struct {
	storage IPremiumStorage
	logger  *my_logger.MyLogger
}

func NewPremiumService(PremiumStorage IPremiumStorage) (*PremiumService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &PremiumService{storage: PremiumStorage, logger: logger}, nil
}

func (p PremiumService) AddPremium(ctx context.Context, productID uint64, userID uint64) error {
	err := p.storage.AddPremium(ctx, productID, userID)
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

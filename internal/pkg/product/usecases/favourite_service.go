package usecases

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/myerrors"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
	"go.uber.org/zap"
)

var _ IFavouriteStorage = (*productrepo.ProductStorage)(nil)

type IFavouriteStorage interface {
	GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error)
	AddToFavourites(ctx context.Context, userID uint64, productID uint64) error
	DeleteFromFavourites(ctx context.Context, userID uint64, productID uint64) error
}

type FavouriteService struct {
	storage IFavouriteStorage
	logger  *zap.SugaredLogger
}

func NewFavouriteService(favouriteStorage IFavouriteStorage) (*FavouriteService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &FavouriteService{storage: favouriteStorage, logger: logger}, nil
}

func (f FavouriteService) GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error) {
	products, err := f.storage.GetUserFavourites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range products {
		product.Sanitize()
	}

	return products, nil
}

func (f FavouriteService) AddToFavourites(ctx context.Context, userID uint64, productID uint64) error {
	err := f.storage.AddToFavourites(ctx, userID, productID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (f FavouriteService) DeleteFromFavourites(ctx context.Context, userID uint64, productID uint64) error {
	err := f.storage.DeleteFromFavourites(ctx, userID, productID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

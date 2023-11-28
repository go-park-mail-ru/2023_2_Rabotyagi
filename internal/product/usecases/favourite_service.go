package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
)

var _ IFavouriteStorage = (*productrepo.ProductStorage)(nil)

type IFavouriteStorage interface {
	GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error)
	AddToFavourites(ctx context.Context, userID uint64, productID uint64) error
	DeleteFromFavourites(ctx context.Context, userID uint64, productID uint64) error
}

type FavouriteService struct {
	storage IFavouriteStorage
	logger  *my_logger.MyLogger
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

func (f FavouriteService) AddToFavourites(ctx context.Context, userID uint64, r io.Reader) error {
	productID := new(models.ProductID)
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(productID); err != nil {
		f.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreOrder)
	}

	err := f.storage.AddToFavourites(ctx, userID, productID.ProductID)
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

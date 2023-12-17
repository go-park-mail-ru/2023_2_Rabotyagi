package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/jackc/pgx/v5"
)

var ErrNoAffectedFavouriteRows = myerrors.NewErrorBadFormatRequest("Не получилось удалить из избранного")

func (p *ProductStorage) selectUserFavourites(ctx context.Context, tx pgx.Tx,
	userID uint64,
) ([]*models.ProductInFeed, error) {
	logger := p.logger.LogReqID(ctx)

	SQLSelectUserFavourites := `SELECT p.id, p.title, p.price, p.city_id,
		p.delivery, p.safe_deal, p.is_active, p.available_count
		FROM public."product" p
		JOIN public."favourite" f ON p.id = f.product_id
		WHERE f.owner_id = $1`

	productsInFavouritesRows, err := tx.Query(ctx, SQLSelectUserFavourites, userID)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curProduct := new(models.ProductInFeed)

	var slProduct []*models.ProductInFeed

	_, err = pgx.ForEachRow(productsInFavouritesRows, []any{
		&curProduct.ID, &curProduct.Title,
		&curProduct.Price, &curProduct.CityID,
		&curProduct.Delivery, &curProduct.SafeDeal, &curProduct.IsActive, &curProduct.AvailableCount,
	}, func() error {
		slProduct = append(slProduct, &models.ProductInFeed{ //nolint:exhaustruct
			ID:             curProduct.ID,
			Title:          curProduct.Title,
			Price:          curProduct.Price,
			CityID:         curProduct.CityID,
			Delivery:       curProduct.Delivery,
			SafeDeal:       curProduct.SafeDeal,
			IsActive:       curProduct.IsActive,
			AvailableCount: curProduct.AvailableCount,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) GetUserFavourites(ctx context.Context, userID uint64) ([]*models.ProductInFeed, error) {
	logger := p.logger.LogReqID(ctx)

	var slProduct []*models.ProductInFeed

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error { //nolint:varnamelen
		slProductInner, err := p.selectUserFavourites(ctx, tx, userID)
		if err != nil {
			return err
		}

		for _, product := range slProductInner {
			productAdditionInner, err := p.getProductAddition(ctx, tx, product.ID, userID)
			if err != nil {
				return err
			}

			product.Images = productAdditionInner.images
			product.Favourites = productAdditionInner.favourites
			product.InFavourites = productAdditionInner.inFavourite

			slProduct = append(slProduct, product)
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slProduct, nil
}

func (p *ProductStorage) addToFavourites(ctx context.Context, tx pgx.Tx, //nolint:varnamelen
	userID uint64, productID uint64,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLAddToFavourites := `INSERT INTO public."favourite"(owner_id, product_id) VALUES($1, $2)`

	_, err := tx.Exec(ctx, SQLAddToFavourites, userID, productID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) AddToFavourites(ctx context.Context, userID uint64, productID uint64) error {
	logger := p.logger.LogReqID(ctx)

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.addToFavourites(ctx, tx, userID, productID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) deleteFromFavourites(ctx context.Context, tx pgx.Tx, //nolint:varnamelen
	userID uint64, productID uint64,
) error {
	logger := p.logger.LogReqID(ctx)

	SQLDeleteFromFavourites := `DELETE FROM public."favourite"
		 WHERE owner_id=$1 AND product_id=$2`

	result, err := tx.Exec(ctx, SQLDeleteFromFavourites, userID, productID)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedFavouriteRows)
	}

	return nil
}

func (p *ProductStorage) DeleteFromFavourites(ctx context.Context, userID uint64, productID uint64) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.deleteFromFavourites(ctx, tx, userID, productID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

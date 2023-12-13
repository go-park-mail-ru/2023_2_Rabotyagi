package repository

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/jackc/pgx/v5"
)

func (p *ProductStorage) addPriceHistoryRecord(ctx context.Context, tx pgx.Tx,
	productID uint64, price uint64) error {
	logger := p.logger.LogReqID(ctx)

	SQLInsertProduct := `INSERT INTO public."price_history"(product_id, price) VALUES($1, $2)`
	_, err := tx.Exec(ctx, SQLInsertProduct, productID, price)

	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) selectPriceHistory(ctx context.Context, tx pgx.Tx,
	productID uint64,
) ([]models.PriceHistoryRecord, error) {
	logger := p.logger.LogReqID(ctx)

	SQLSelectPriceHistory := `SELECT price, created_at
		FROM public."price_history" 
		WHERE product_id = $1
		ORDER BY created_at ASC`

	productsInFavouritesRows, err := tx.Query(ctx, SQLSelectPriceHistory, productID)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curRecord := new(models.PriceHistoryRecord)

	var slRecord []models.PriceHistoryRecord

	_, err = pgx.ForEachRow(productsInFavouritesRows, []any{
		&curRecord.Price, &curRecord.CreatedAt,
	}, func() error {
		slRecord = append(slRecord, models.PriceHistoryRecord{ //nolint:exhaustruct
			Price:     curRecord.Price,
			CreatedAt: curRecord.CreatedAt,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slRecord, nil
}

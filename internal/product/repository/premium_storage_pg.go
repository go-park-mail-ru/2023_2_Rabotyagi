package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/jackc/pgx/v5"
)

var (
	ErrPremiumStatusNotFound = myerrors.NewErrorBadFormatRequest(
		"Не найдено объявление с таким id у такого пользователя")
	ErrNoAffectedPremiumStatusRows = myerrors.NewErrorBadFormatRequest(
		"Не получилось обновить статус премиума объявления")
)

func (p *ProductStorage) UpdateStatusPremium(ctx context.Context, status uint8, productID uint64, userID uint64) error {
	SQLAddPremium := `UPDATE public."product" 
SET premium_status=$1 WHERE id=$4 AND saler_id=$5`

	result, err := p.pool.Exec(ctx, SQLAddPremium, statuses.IntStatusPremiumSucceeded, productID, userID)
	if err != nil {
		p.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		err = fmt.Errorf("%w productID=%d userId=%d status=%d", ErrNoAffectedPremiumStatusRows, productID, userID, status)
		p.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) addPremium(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64,
	premiumBegin time.Time, premiumExpire time.Time,
) error {
	SQLAddPremium := `UPDATE public."product" 
SET premium_status=$1, premium_begin=$2, premium_expire=$3 WHERE id=$4 AND saler_id=$5`

	result, err := tx.Exec(ctx, SQLAddPremium, statuses.IntStatusPremiumSucceeded,
		premiumBegin, premiumExpire, productID, userID)
	if err != nil {
		p.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedProductRows)
	}

	return nil
}

func (p *ProductStorage) AddPremium(ctx context.Context, productID uint64, userID uint64,
	premiumBegin time.Time, premiumExpire time.Time,
) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.addPremium(ctx, tx, productID, userID, premiumBegin, premiumExpire)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		p.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductStorage) selectPremiumStatusOfProduct(ctx context.Context,
	productID uint64, userID uint64,
) (uint8, error) {
	logger := p.logger.LogReqID(ctx)

	var premiumStatus uint8

	SQLSelectPremiumStatus := `SELECT premium_status FROM public."product" WHERE id=$1 AND saler_id=$2`

	productRow := p.pool.QueryRow(ctx, SQLSelectPremiumStatus, productID, userID)
	if err := productRow.Scan(&premiumStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf(myerrors.ErrTemplate, ErrPremiumStatusNotFound)
		}

		logger.Errorf("error with productId=%d and userID=%d: %+v", productID, userID, err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return premiumStatus, nil
}

func (p *ProductStorage) CheckPremiumStatus(ctx context.Context, productID uint64, userID uint64) (uint8, error) {
	status, err := p.selectPremiumStatusOfProduct(ctx, productID, userID)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return status, nil
}

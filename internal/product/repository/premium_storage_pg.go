package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/jackc/pgx/v5"
)

func (p *ProductStorage) addPremium(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64,
	premiumBegin time.Time, premiumExpire time.Time) error {
	SQLAddPremium := `UPDATE public."product" SET premium=true, premium_begin=$1, premium_expire=$2 WHERE id=$3 AND saler_id=$4`

	result, err := tx.Exec(ctx, SQLAddPremium, premiumBegin, premiumExpire, productID, userID)
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
	premiumBegin time.Time, premiumExpire time.Time) error {
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

func (p *ProductStorage) removePremium(ctx context.Context, tx pgx.Tx, productID uint64, userID uint64) error {
	SQLRemovePremium := `UPDATE public."product" SET premium=false WHERE id=$1 AND saler_id=$2`

	result, err := tx.Exec(ctx, SQLRemovePremium, productID, userID)
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

func (p *ProductStorage) RemovePremium(ctx context.Context, productID uint64, userID uint64) error {
	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		err := p.removePremium(ctx, tx, productID, userID)
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

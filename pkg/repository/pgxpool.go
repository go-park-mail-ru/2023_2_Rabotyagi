package repository

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, urlDataBase string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, urlDataBase)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return pool, nil
}

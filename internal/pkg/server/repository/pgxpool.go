package repository

import (
	"context"
	"fmt"
	"log"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, urlDataBase string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, urlDataBase)
	if err != nil {
		log.Printf("Error init db connection: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Error ping db: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return pool, nil
}

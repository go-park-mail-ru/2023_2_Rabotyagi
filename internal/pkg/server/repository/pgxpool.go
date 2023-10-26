package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, URLDataBase string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, URLDataBase)
	if err != nil {
		log.Printf("Error init db connection: %v\n", err)

		return nil, err //nolint:wrapcheck
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Error ping db: %v\n", err)

		return nil, err //nolint:wrapcheck
	}

	return pool, nil
}

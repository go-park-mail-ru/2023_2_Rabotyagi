package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type CategoryStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewCategoryStorage(pool *pgxpool.Pool, logger *zap.SugaredLogger) *CategoryStorage {
	return &CategoryStorage{
		pool:   pool,
		logger: logger,
	}
}

func (c *CategoryStorage) selectFullCatgories(ctx context.Context, tx pgx.Tx) ([]*models.Category, error) {
	var categories []*models.Category

	SQLSelectFullCatgories := `SELECT "category".id,"category".name, "category".parent_id FROM public."category"`

	categoriesRows, err := tx.Query(ctx, SQLSelectFullCatgories)
	if err != nil {
		c.logger.Errorf("in selectOrdersInBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curCat := new(models.Category)

	_, err = pgx.ForEachRow(categoriesRows, []any{
		&curCat.ID, &curCat.Name, &curCat.ParentID,
	}, func() error {
		categories = append(categories, &models.Category{ //nolint:exhaustruct
			ID:       curCat.ID,
			Name:     curCat.Name,
			ParentID: curCat.ParentID,
		})

		return nil
	})
	if err != nil {
		c.logger.Errorf("in selectOrdersInBasketByUserID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return categories, nil
}

func (c *CategoryStorage) GetFullCategories(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category

	err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		categoriesInner, err := c.selectFullCatgories(ctx, tx)
		if err != nil {
			return err
		}

		categories = categoriesInner

		return nil
	})
	if err != nil {
		c.logger.Errorf("in GetFullCatgories: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return categories, nil
}

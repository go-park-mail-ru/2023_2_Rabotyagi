package repository

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/myerrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type CategoryStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewCategoryStorage(pool *pgxpool.Pool) (*CategoryStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &CategoryStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (c *CategoryStorage) selectFullCatgories(ctx context.Context, tx pgx.Tx) ([]*models.Category, error) {
	var categories []*models.Category

	SQLSelectFullCatgories := `SELECT "category".id,"category".name, "category".parent_id FROM public."category"`

	categoriesRows, err := tx.Query(ctx, SQLSelectFullCatgories)
	if err != nil {
		c.logger.Errorln(err)

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
		c.logger.Errorln(err)

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
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return categories, nil
}

package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/pgxpool"
	"strings"

	"github.com/jackc/pgx/v5"
)

type CategoryStorage struct {
	pool   pgxpool.IPgxPool
	logger *my_logger.MyLogger
}

func NewCategoryStorage(pool pgxpool.IPgxPool) (*CategoryStorage, error) {
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
	logger := c.logger.LogReqID(ctx)

	var categories []*models.Category

	SQLSelectFullCatgories := `SELECT "category".id,"category".name, "category".parent_id FROM public."category"`

	categoriesRows, err := tx.Query(ctx, SQLSelectFullCatgories)
	if err != nil {
		logger.Errorln(err)

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
		logger.Errorln(err)

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

func (c *CategoryStorage) searchCategory(ctx context.Context, tx pgx.Tx, searchInput string) ([]*models.Category, error) {
	logger := c.logger.LogReqID(ctx)

	SQLSearchCategory := `SELECT category.id, category.name, category.parent_id
						FROM public."category"
						WHERE LOWER(name) LIKE $1 
						LIMIT 5;`

	var cities []*models.Category

	categoriesRows, err := tx.Query(ctx, SQLSearchCategory, "%"+strings.ToLower(searchInput)+"%")
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curCategory := new(models.Category)

	_, err = pgx.ForEachRow(categoriesRows, []any{
		&curCategory.ID, &curCategory.Name, &curCategory.ParentID,
	}, func() error {
		cities = append(cities, &models.Category{ //nolint:exhaustruct
			ID:       curCategory.ID,
			Name:     curCategory.Name,
			ParentID: curCategory.ParentID,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return cities, nil
}

func (c *CategoryStorage) SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error) {
	logger := c.logger.LogReqID(ctx)

	var categories []*models.Category

	err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		categoriesInner, err := c.searchCategory(ctx, tx, searchInput)
		if err != nil {
			return err
		}

		categories = categoriesInner

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return categories, nil
}

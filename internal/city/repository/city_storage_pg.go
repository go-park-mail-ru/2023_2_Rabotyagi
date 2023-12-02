package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"strings"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/pgxpool"
	"github.com/jackc/pgx/v5"
)

type CityStorage struct {
	pool   pgxpool.IPgxPool
	logger *my_logger.MyLogger
}

func NewCityStorage(pool pgxpool.IPgxPool) (*CityStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &CityStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (c *CityStorage) selectFullCities(ctx context.Context, tx pgx.Tx) ([]*models.City, error) {
	logger := c.logger.LogReqID(ctx)

	var cities []*models.City

	SQLSelectFullCities := `SELECT "city".id,"city".name FROM public."city"`

	citiesRows, err := tx.Query(ctx, SQLSelectFullCities)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curCity := new(models.City)

	_, err = pgx.ForEachRow(citiesRows, []any{
		&curCity.ID, &curCity.Name,
	}, func() error {
		cities = append(cities, &models.City{ //nolint:exhaustruct
			ID:   curCity.ID,
			Name: curCity.Name,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return cities, nil
}

func (c *CityStorage) GetFullCities(ctx context.Context) ([]*models.City, error) {
	var cities []*models.City

	err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		citiesInner, err := c.selectFullCities(ctx, tx)
		if err != nil {
			return err
		}

		cities = citiesInner

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return cities, nil
}

func (c *CityStorage) searchCity(ctx context.Context, tx pgx.Tx, searchInput string) ([]*models.City, error) {
	logger := c.logger.LogReqID(ctx)

	SQLSearchCity := `SELECT city.id, city.name
						FROM public."city"
						WHERE LOWER(name) LIKE $1 
						LIMIT 5;`

	var cities []*models.City

	citiesRows, err := tx.Query(ctx, SQLSearchCity, "%"+strings.ToLower(searchInput)+"%")
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curCity := new(models.City)

	_, err = pgx.ForEachRow(citiesRows, []any{
		&curCity.ID, &curCity.Name,
	}, func() error {
		cities = append(cities, &models.City{ //nolint:exhaustruct
			ID:   curCity.ID,
			Name: curCity.Name,
		})

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return cities, nil
}

func (c *CityStorage) SearchCity(ctx context.Context, searchInput string) ([]*models.City, error) {
	logger := c.logger.LogReqID(ctx)

	var cities []*models.City

	err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		citiesInner, err := c.searchCity(ctx, tx, searchInput)
		if err != nil {
			return err
		}

		cities = citiesInner

		return nil
	})
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return cities, nil
}

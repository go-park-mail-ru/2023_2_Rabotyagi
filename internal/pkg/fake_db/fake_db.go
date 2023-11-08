package fake_db

import (
	"context"
	"github.com/jackc/pgx/v5"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/fake_db/repository"
	serverrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/repository"

	"go.uber.org/zap"
)

func RunScriptFillDB(URLDataBase string, logger *zap.SugaredLogger, baseCount uint) error {

	userMaxCount := baseCount
	categoryMaxCount := userMaxCount/10 + 1
	productMaxCount := userMaxCount * 4
	orderMaxCount := userMaxCount * 2
	favouritesMaxCount := userMaxCount
	baseCtx := context.Background()

	pool, err := serverrepo.NewPgxPool(baseCtx, URLDataBase)
	if err != nil {
		logger.Error(err)

		return err
	}

	fakeStorage := repository.FakeStorage{Pool: pool, Logger: logger}

	err = pgx.BeginFunc(baseCtx, pool, func(tx pgx.Tx) error {
		err = fakeStorage.InsertUsersWithoutID(baseCtx, tx, userMaxCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err)

		return err
	}

	err = pgx.BeginFunc(baseCtx, pool, func(tx pgx.Tx) error {
		err = fakeStorage.InsertCategories(baseCtx, tx, categoryMaxCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err)

		return err
	}

	err = pgx.BeginFunc(baseCtx, pool, func(tx pgx.Tx) error {
		err = fakeStorage.InsertProducts(baseCtx,
			tx, productMaxCount, userMaxCount, categoryMaxCount,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err)

		return err
	}

	err = pgx.BeginFunc(baseCtx, pool, func(tx pgx.Tx) error {
		err = fakeStorage.InsertOrders(baseCtx,
			tx, userMaxCount, orderMaxCount, productMaxCount,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err)

		return err
	}

	err = pgx.BeginFunc(baseCtx, pool, func(tx pgx.Tx) error {
		err = fakeStorage.InsertFavourites(baseCtx,
			tx, favouritesMaxCount, userMaxCount, productMaxCount,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

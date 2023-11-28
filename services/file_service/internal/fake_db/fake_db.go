package fake_db

import (
	"context"
	serverrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/fake_db/repository"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func RunScriptFillDB(URLDataBase string,
	logger *zap.SugaredLogger, baseCount uint, pathToRoot string,
) error {
	prefixURL := "img/"
	maxNameImg := uint(12)
	userMaxCount := baseCount
	categoryMaxCount := userMaxCount/10 + 1
	cityMaxCount := userMaxCount/10 + 1
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
		err = fakeStorage.InsertCity(baseCtx, tx, cityMaxCount)
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
			tx, productMaxCount, userMaxCount, categoryMaxCount, cityMaxCount,
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
		err = fakeStorage.InsertImages(baseCtx,
			tx, maxNameImg, productMaxCount, prefixURL, pathToRoot,
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

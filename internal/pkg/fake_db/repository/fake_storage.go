package repository

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/fake_db/usecases"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FakeStorage struct {
	Pool   *pgxpool.Pool
	Logger *zap.SugaredLogger
}

func (f *FakeStorage) InsertUsersWithoutID(ctx context.Context, tx pgx.Tx, count uint) error {
	slUser := [][]any{}
	columns := []string{"email", "phone", "name", "password", "birthday"}

	f.Logger.Infof("start filling users")

	for i := 0; i < int(count); i++ {
		if i%(int(count)%100) == 0 {
			f.Logger.Infof("filled i=%d of %d users", i, count)
		}

		user, err := usecases.FakeUserWihtoutID(i)
		if err != nil {
			f.Logger.Error(err)

			return err
		}

		slUser = append(slUser, []any{user.Email, user.Phone, user.Name, user.Password, user.Birthday})
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "user"},
		columns,
		pgx.CopyFromRows(slUser),
	)

	if err != nil {
		f.Logger.Error(err)

		return err
	}

	f.Logger.Infof("end filling users\n")

	return nil
}

// InsertCategories open new connection because categories have constraint referenses on parent_id.
// At this reason I insert parent categories in second connection
func (f *FakeStorage) InsertCategories(ctx context.Context, tx pgx.Tx, count uint) error {
	count %= 10
	count++

	slBaseCategories := [][]any{}
	slCategories := [][]any{}
	baseColumns := []string{"name"}
	columns := []string{"name", "parent_id"}

	f.Logger.Infof("start filling categories")

	categories := gofakeit.Categories()
	idxCategory := 1
	idxTotal := 1

	for key, subCategory := range categories {
		if idxCategory > int(count) {
			break
		}

		slBaseCategories = append(slBaseCategories, []any{key})
		idxTotal++

		for _, nameSubCategory := range subCategory {
			slCategories = append(slCategories, []any{nameSubCategory, idxCategory})
			idxTotal++
		}

		idxCategory += idxTotal
	}

	initStat := f.Pool.Stat()

	countCopyBase, err := f.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"public", "category"},
		baseColumns,
		pgx.CopyFromRows(slBaseCategories),
	)
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	beforeStat := f.Pool.Stat()

	for i := 0; ; i++ {
		if initStat.AcquiredConns() == beforeStat.AcquiredConns() {
			break
		}

		time.Sleep(time.Millisecond * 100)
		i++

		if i > 10 {
			return fmt.Errorf("Не дождались возврата коннекта для записи родительских категория")
		}
	}

	countCopy, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "category"},
		columns,
		pgx.CopyFromRows(slCategories),
	)
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	if countCopyBase == 0 || countCopy == 0 {
		return fmt.Errorf("countCopyBase=%d and coutCopy=%d", countCopyBase, countCopy)
	}

	f.Logger.Infof("end filling %d categories", countCopy+countCopyBase)

	return nil
}

func (f *FakeStorage) InsertProducts(ctx context.Context, tx pgx.Tx, count uint) error {
	productCount := 2 * count
	slProduct := [][]any{}
	columns := []string{
		"saler_id", "category_id", "title", "description", "price",
		"available_count", "city", "delivery", "safe_deal",
	}

	f.Logger.Infof("start filling users")

	for i := 0; i < int(productCount); i++ {
		if i%(int(count)%100) == 0 {
			f.Logger.Infof("filled i=%d of %d products", i, productCount)
		}

		preProduct := usecases.FakePreProduct(count)

		slProduct = append(slProduct,
			[]any{
				preProduct.SalerID, preProduct.CategoryID, preProduct.Title,
				preProduct.Description, preProduct.Price, preProduct.AvailableCount, preProduct.City,
				preProduct.Delivery, preProduct.SafeDeal,
			},
		)
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "product"},
		columns,
		pgx.CopyFromRows(slProduct),
	)
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	f.Logger.Infof("end filling products\n")

	return nil

}

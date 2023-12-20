package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/fake_db/usecases"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FakeStorage struct {
	Pool   *pgxpool.Pool
	Logger *mylogger.MyLogger
}

func (f *FakeStorage) InsertUsersWithoutID(ctx context.Context, tx pgx.Tx, userCount uint) error {
	slUser := [][]any{}
	columns := []string{"email", "phone", "name", "password", "birthday"}

	f.Logger.Infof("start filling users")

	for i := 0; i < int(userCount); i++ {
		if i%(int(userCount)/100+1) == 0 {
			f.Logger.Infof("filled i=%d of %d users", i, userCount)
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

func (f *FakeStorage) InsertCity(ctx context.Context, tx pgx.Tx, cityCount uint) error {
	slCity := [][]any{}
	columns := []string{"name"}

	f.Logger.Infof("start filling users")

	for i := 0; i < int(cityCount); i++ {
		if i%(int(cityCount)/100+1) == 0 {
			f.Logger.Infof("filled i=%d of %d city", i, cityCount)
		}

		city := usecases.FakeCity()

		slCity = append(slCity, []any{city})
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "city"},
		columns,
		pgx.CopyFromRows(slCity),
	)

	if err != nil {
		f.Logger.Error(err)

		return err
	}

	f.Logger.Infof("end filling city\n")

	return nil
}

// InsertCategories open new connection because categories have constraint referenses on parent_id.
// At this reason I insert parent categories in second connection
func (f *FakeStorage) InsertCategories(ctx context.Context, tx pgx.Tx, categoriesCount uint) error {
	categoriesCount++

	slBaseCategories := [][]any{}
	slCategories := [][]any{}
	baseColumns := []string{"name"}
	columns := []string{"name", "parent_id"}

	f.Logger.Infof("start filling categories")

	categories := gofakeit.Categories()
	idxCategory := 1
	idxTotal := 1

	for key, subCategory := range categories {
		if idxCategory > int(categoriesCount) {
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

func (f *FakeStorage) InsertProducts(ctx context.Context,
	tx pgx.Tx, productCount uint, userMaxCount uint, categoryMaxCount uint, cityMaxCount uint,
) error {
	slProduct := [][]any{}
	columns := []string{
		"saler_id", "category_id", "city_id", "title", "description", "price",
		"available_count", "delivery", "safe_deal", "views", "premium",
	}

	f.Logger.Infof("start filling users")

	for i := 0; i < int(productCount); i++ {
		if i%(int(productCount)/100+1) == 0 {
			f.Logger.Infof("filled i=%d of %d products", i, productCount)
		}

		preProduct := usecases.FakeProduct(userMaxCount, categoryMaxCount, cityMaxCount)

		slProduct = append(slProduct,
			[]any{
				preProduct.SalerID, preProduct.CategoryID, preProduct.CityID, preProduct.Title,
				preProduct.Description, preProduct.Price, preProduct.AvailableCount,
				preProduct.Delivery, preProduct.SafeDeal, preProduct.Views, preProduct.Premium,
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

func (f *FakeStorage) InsertOrders(ctx context.Context,
	tx pgx.Tx, userMaxCount uint, ordersMaxCount uint, productMaxCount uint,
) error {
	slOrder := [][]any{}
	columns := []string{
		"owner_id", "product_id", "count",
	}

	f.Logger.Infof("start filling orders")

	for i := 0; i < int(ordersMaxCount); i++ {
		if i%(int(ordersMaxCount)/100+1) == 0 {
			f.Logger.Infof("filled i=%d of %d orders", i, ordersMaxCount)
		}

		preOrder := usecases.FakePreOrder(userMaxCount, productMaxCount)

		slOrder = append(slOrder,
			[]any{preOrder.OwnerID, preOrder.ProductID, preOrder.Count},
		)
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "order"},
		columns,
		pgx.CopyFromRows(slOrder),
	)
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	f.Logger.Infof("end filling orders\n")

	return nil
}

func (f *FakeStorage) InsertImages(ctx context.Context,
	tx pgx.Tx, maxNameImage uint, maxCountProducts uint, prefixURL string, pathToRoot string,
) error {
	fakeGeneratorImg, err := usecases.NewFakeGeneratorImg(maxNameImage, prefixURL, pathToRoot)
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	slImg := [][]any{}
	columns := []string{
		"url", "product_id",
	}

	f.Logger.Infof("start filling images")

	for i := 1; i < int(maxCountProducts); i++ {
		if i%(int(maxCountProducts)%100) == 0 {
			f.Logger.Infof("filled images i=%d of %d prodcuts", i, maxCountProducts)
		}

		URLs := fakeGeneratorImg.GetURLs(uint(gofakeit.Number(0, int(maxNameImage))))

		for _, url := range URLs {
			slImg = append(slImg, []any{url, i})
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "image"},
		columns,
		pgx.CopyFromRows(slImg),
	)
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	f.Logger.Infof("end filling images\n")

	return nil
}

// InsertFavourites TODO fix troubles with uniq together
func (f *FakeStorage) InsertFavourites(ctx context.Context,
	tx pgx.Tx, maxCountFavourites uint, maxCountUsers uint, maxCountProducts uint,
) error {
	slOrder := [][]any{}
	columns := []string{
		"owner_id", "product_id",
	}

	f.Logger.Infof("start filling favourites")

	for i := 0; i < int(maxCountFavourites); i++ {
		if i%(int(maxCountFavourites)/100+1) == 0 {
			f.Logger.Infof("filled i=%d of %d favourites", i, maxCountFavourites)
		}

		ownerID, productID := usecases.FakeFavourite(maxCountUsers, maxCountProducts)

		slOrder = append(slOrder,
			[]any{ownerID, productID},
		)
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"public", "favourite"},
		columns,
		pgx.CopyFromRows(slOrder))
	if err != nil {
		f.Logger.Error(err)

		return err
	}

	f.Logger.Infof("end filling favourites\n")

	return nil
}

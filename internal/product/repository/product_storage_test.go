package repository_test

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/pashagolub/pgxmock/v3"
	"testing"
	"time"
)

func TestCloseProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		salerID                uint64
		productID              uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			salerID:   1,
			productID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage)

			err = catStorage.CloseProduct(ctx, testCase.productID, testCase.salerID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestActivateProduct(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		salerID                uint64
		productID              uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			salerID:   1,
			productID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage)

			err = catStorage.ActivateProduct(ctx, testCase.productID, testCase.salerID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		salerID                uint64
		productID              uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			salerID:   1,
			productID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage)

			err = catStorage.DeleteProduct(ctx, testCase.productID, testCase.salerID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSearchProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		searchInput            string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT title FROM product`).WithArgs("Ca").
					WillReturnRows(pgxmock.NewRows([]string{"title"}).
						AddRow("Car"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput: "Ca",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage)

			_, err = catStorage.SearchProduct(ctx, testCase.searchInput)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetSearchProductFeed(t *testing.T) { //nolint:funlen
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		searchInput            string
		lastNumber             uint64
		limit                  uint64
		userID                 uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT id, title, price, city_id, 
       delivery, safe_deal, is_active, available_count FROM product`).WithArgs("Ca", uint64(0), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "title", "price", "city_id", //nolint:gofumpt
						"delivery", "safe_deal", "is_active", "available_count"}).
						AddRow(uint64(1), "Car", uint64(1212), uint64(6), true, true, true, uint32(2)))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput: "Ca",
			lastNumber:  0,
			limit:       1,
			userID:      1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			prodStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(prodStorage)

			_, err = prodStorage.GetSearchProductFeed(ctx, testCase.searchInput,
				testCase.lastNumber, testCase.limit, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		productID              uint64
		userID                 uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city_id,
       delivery, safe_deal, is_active, premium FROM public."product" `).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"saler_id", "category_id", "title", "description", "price", //nolint:gofumpt
						"created_at", "views", "available_count", "city_id", "delivery",
						"safe_deal", "is_active", "premium"}).
						AddRow(uint64(2), uint64(1), "Car", "text", uint64(1212), time.Now(),
							uint32(6), uint32(4), uint64(6), true, true, true, false))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectQuery(`SELECT price, created_at FROM public."price_history"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"price", "created_at"}).
						AddRow(uint64(123123), time.Now()))

				mockPool.ExpectQuery(`SELECT EXISTS`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).
						AddRow(false))

				mockPool.ExpectExec(`INSERT INTO public."view"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:    1,
			productID: 1,
		},
		{
			name: "test my",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city_id,
       delivery, safe_deal, is_active, premium FROM public."product" `).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"saler_id", "category_id", "title", "description", "price", //nolint:gofumpt
						"created_at", "views", "available_count", "city_id", "delivery",
						"safe_deal", "is_active", "premium"}).
						AddRow(uint64(1), uint64(1), "Car", "text", uint64(1212), time.Now(),
							uint32(6), uint32(4), uint64(6), true, true, true, true))

				mockPool.ExpectQuery(`SELECT premium_expire FROM public."product"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"premium_expire"}).
						AddRow(sql.NullTime{Valid: true, Time: time.Now()}))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectQuery(`SELECT price, created_at FROM public."price_history"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"price", "created_at"}).
						AddRow(uint64(123123), time.Now()))

				mockPool.ExpectQuery(`SELECT EXISTS`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).
						AddRow(false))

				mockPool.ExpectExec(`INSERT INTO public."view"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:    1,
			productID: 1,
		},
	}

	for _, testCase := range testCases { //paralleltest
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			prodStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(prodStorage)

			_, err = prodStorage.GetProduct(ctx, testCase.productID, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

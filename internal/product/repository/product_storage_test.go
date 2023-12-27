package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/pashagolub/pgxmock/v3"
)

func TestCloseProduct(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		salerID                uint64
		productID              uint64
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			salerID:       1,
			productID:     1,
			expectedError: nil,
		},
		{
			name: "test no rows affected",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			salerID:       1,
			productID:     1,
			expectedError: repository.ErrNoAffectedProductRows,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage, mockPool)

			errActual := catStorage.CloseProduct(ctx, testCase.productID, testCase.salerID)

			err = utils.EqualError(errActual, testCase.expectedError)
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

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		salerID                uint64
		productID              uint64
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			salerID:       1,
			productID:     1,
			expectedError: nil,
		},
		{
			name: "test no affected rows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			salerID:       1,
			productID:     1,
			expectedError: repository.ErrNoAffectedProductRows,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage, mockPool)

			errActual := catStorage.ActivateProduct(ctx, testCase.productID, testCase.salerID)

			err = utils.EqualError(errActual, testCase.expectedError)
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

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		salerID                uint64
		productID              uint64
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			salerID:       1,
			productID:     1,
			expectedError: nil,
		},
		{
			name: "test no affected rows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."product"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			salerID:       1,
			productID:     1,
			expectedError: repository.ErrNoAffectedProductRows,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage, mockPool)

			errActual := catStorage.DeleteProduct(ctx, testCase.productID, testCase.salerID)

			err = utils.EqualError(errActual, testCase.expectedError)
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

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		searchInput            string
		expectedResponse       []string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT title`).WithArgs("Ca").
					WillReturnRows(pgxmock.NewRows([]string{"title"}).
						AddRow("Car"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput:      "Ca",
			expectedResponse: []string{"Car"},
		},
		{
			name: "test more rows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT title`).WithArgs("Ca").
					WillReturnRows(pgxmock.NewRows([]string{"title"}).
						AddRow("Car").AddRow("Cat").AddRow("Carrot"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput:      "Ca",
			expectedResponse: []string{"Car", "Cat", "Carrot"},
		},
		{
			name: "test empty",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT title`).WithArgs("Ca").
					WillReturnRows(pgxmock.NewRows([]string{}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput:      "Ca",
			expectedResponse: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage, mockPool)

			response, err := catStorage.SearchProduct(ctx, testCase.searchInput)
			if err != nil {
				t.Fatal(err)
			}

			if err := utils.EqualTest(response, testCase.expectedResponse); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetSearchProductFeed(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		searchInput            string
		lastNumber             uint64
		limit                  uint64
		userID                 uint64
		expectedResponse       []*models.ProductInFeed
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT id, title, price, city_id, 
       delivery, safe_deal, is_active, available_count, premium_status FROM product`).
					WithArgs("Ca", uint64(0), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "title", "price", "city_id",
						"delivery", "safe_deal", "is_active", "available_count", "premium",
					}).
						AddRow(uint64(1), "Car", uint64(1212), uint64(6), true, true, true, uint32(2), statuses.IntStatusPremiumNot))

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
			expectedResponse: []*models.ProductInFeed{
				{
					ID: 1, Title: "Car", Price: 1212, CityID: 6, Delivery: true, IsActive: true,
					SafeDeal: true, AvailableCount: 2, Premium: false,
					Images: []models.Image{{URL: "safsafddasf"}}, InFavourites: true, Favourites: 1,
				},
			},
		},
		{
			name: "test more data",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT id, title, price, city_id, 
       delivery, safe_deal, is_active, available_count, premium_status FROM product`).
					WithArgs("Ca", uint64(0), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "title", "price", "city_id",
						"delivery", "safe_deal", "is_active", "available_count", "premium",
					}).
						AddRow(uint64(1), "Car", uint64(1212), uint64(6), true, true, true, uint32(2), statuses.IntStatusPremiumNot).
						AddRow(uint64(2), "Cat", uint64(1212), uint64(6), true, true, true, uint32(2), statuses.IntStatusPremiumNot).
						AddRow(uint64(3), "Carrot", uint64(1212), uint64(6), true, true, true, uint32(2), statuses.IntStatusPremiumNot))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(2)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(2)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(2), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(3)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(3)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(3), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput: "Ca",
			lastNumber:  0,
			limit:       1,
			userID:      1,
			expectedResponse: []*models.ProductInFeed{
				{
					ID: 1, Title: "Car", Price: 1212, CityID: 6, Delivery: true, IsActive: true, SafeDeal: true,
					AvailableCount: 2, Premium: false, Images: []models.Image{{URL: "safsafddasf"}}, InFavourites: true, Favourites: 1,
				},
				{
					ID: 2, Title: "Cat", Price: 1212, CityID: 6, Delivery: true, IsActive: true, SafeDeal: true,
					AvailableCount: 2, Premium: false, Images: []models.Image{{URL: "safsafddasf"}}, InFavourites: true, Favourites: 1,
				},
				{
					ID: 3, Title: "Carrot", Price: 1212, CityID: 6, Delivery: true, IsActive: true, SafeDeal: true,
					AvailableCount: 2, Premium: false, Images: []models.Image{{URL: "safsafddasf"}}, InFavourites: true, Favourites: 1,
				},
			},
		},
		{
			name: "test empty",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT id, title, price, city_id, 
       delivery, safe_deal, is_active, available_count, premium_status FROM product`).
					WithArgs("Ca", uint64(0), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			searchInput:      "Ca",
			lastNumber:       0,
			limit:            1,
			userID:           1,
			expectedResponse: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			prodStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(prodStorage, mockPool)

			response, err := prodStorage.GetSearchProductFeed(ctx, testCase.searchInput,
				testCase.lastNumber, testCase.limit, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := utils.EqualTest(response, testCase.expectedResponse); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		productID              uint64
		userID                 uint64
		expectedResponse       *models.Product
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city_id,
       delivery, safe_deal, is_active, premium_status FROM public."product" `).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{
						"saler_id", "category_id", "title",
						"description", "price", "created_at", "views", "available_count", "city_id",
						"delivery", "safe_deal", "is_active", "premium_status",
					}).
						AddRow(uint64(2), uint64(1), "Car", "text", uint64(1212), time.Time{},
							uint32(6), uint32(4), uint64(6), true, true, true, statuses.IntStatusPremiumNot))

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
						AddRow(uint64(123123), time.Time{}))

				mockPool.ExpectQuery(`SELECT id FROM public."comment" WHERE sender_id=\$1 AND recipient_id=\$2`).
					WithArgs(uint64(1), uint64(2)).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint64(1)))

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
			expectedResponse: &models.Product{ //nolint:exhaustruct
				ID: 1, SalerID: 2, CategoryID: 1, CityID: 6, Title: "Car", Description: "text", Price: 1212, CreatedAt: time.Time{},
				Views: 6, AvailableCount: 4, Delivery: true, SafeDeal: true, InFavourites: true, IsActive: true, Premium: false,
				Images: []models.Image{{URL: "safsafddasf"}},
				PriceHistory: []models.PriceHistoryRecord{{
					Price:     123123,
					CreatedAt: time.Time{},
				}}, Favourites: 1, CommentID: sql.NullInt64{Valid: true, Int64: 1},
			},
		},
		{
			name: "test my",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city_id,
       delivery, safe_deal, is_active, premium_status FROM public."product" `).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{
						"saler_id", "category_id", "title",
						"description", "price", "created_at", "views", "available_count", "city_id",
						"delivery", "safe_deal", "is_active", "premium_status",
					}).
						AddRow(uint64(1), uint64(1), "Car", "text", uint64(1212), time.Time{},
							uint32(6), uint32(4), uint64(6), true, true, true, statuses.IntStatusPremiumNot))

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
						AddRow(uint64(123123), time.Time{}))

				mockPool.ExpectQuery(`SELECT id FROM public."comment" WHERE sender_id=\$1 AND recipient_id=\$2`).
					WithArgs(uint64(1), uint64(1)).WillReturnRows(pgxmock.NewRows([]string{"id"}))

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
			expectedResponse: &models.Product{ //nolint:exhaustruct
				ID: 1, SalerID: 1, CategoryID: 1, CityID: 6, Title: "Car", Description: "text", Price: 1212, CreatedAt: time.Time{},
				Views: 6, AvailableCount: 4, Delivery: true, SafeDeal: true, InFavourites: true, IsActive: true, Premium: false,
				Images:       []models.Image{{URL: "safsafddasf"}},
				PriceHistory: []models.PriceHistoryRecord{{Price: 123123, CreatedAt: time.Time{}}}, Favourites: 1,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			ctx := context.Background()

			prodStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(prodStorage, mockPool)

			response, err := prodStorage.GetProduct(ctx, testCase.productID, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := utils.EqualTest(response, testCase.expectedResponse); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

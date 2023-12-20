package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/pashagolub/pgxmock/v3"
)

func TestDeleteOrder(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
		orderID                uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."order"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:  1,
			orderID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			basketStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(basketStorage)

			err = basketStorage.DeleteOrder(ctx, testCase.orderID, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateOrderCount(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		ownerID                uint64
		orderID                uint64
		newCount               uint32
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."order"`).WithArgs(uint32(4), uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:  1,
			orderID:  1,
			newCount: 4,
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

			err = catStorage.UpdateOrderCount(ctx, testCase.ownerID, testCase.orderID, testCase.newCount)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		ownerID                uint64
		orderID                uint64
		newStatus              uint8
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT status, count FROM public."order"`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"status", "count"}).
						AddRow(uint8(1), uint32(1)))

				mockPool.ExpectExec(`UPDATE public."order"`).WithArgs(uint8(2), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:   1,
			orderID:   1,
			newStatus: 2,
		},
		{
			name: "test status default not 0",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT status, count FROM public."order"`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"status", "count"}).
						AddRow(uint8(0), uint32(1)))

				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(uint32(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectExec(`UPDATE public."order"`).WithArgs(uint8(2), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:   1,
			orderID:   1,
			newStatus: 2,
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

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage, mockPool)

			err = catStorage.UpdateOrderStatus(ctx, testCase.ownerID, testCase.orderID, testCase.newStatus)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddOrderInBasket(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
		productID              uint64
		count                  uint32
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT saler_id, category_id, title,
       description, price, created_at, views, available_count, city_id,
       delivery, safe_deal, is_active, premium FROM public."product" `).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"saler_id", "category_id", "title", "description", "price",
						"created_at", "views", "available_count", "city_id", "delivery",
						"safe_deal", "is_active", "premium"}).
						AddRow(uint64(1), uint64(1), "Car", "text", uint64(1212), time.Now(),
							uint32(6), uint32(4), uint64(6), true, true, true, false))

				mockPool.ExpectExec(`INSERT INTO public."order"`).WithArgs(uint64(1), uint64(1), uint32(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectQuery(`SELECT last_value FROM "public"."order_id_seq";`).
					WillReturnRows(pgxmock.NewRows([]string{"last_value"}).
						AddRow(uint64(1)))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:    1,
			productID: 1,
			count:     1,
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

			_, err = catStorage.AddOrderInBasket(ctx, testCase.userID, testCase.productID, testCase.count)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetOrdersInBasket(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city_id, "order".count, "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM public."order"
    INNER JOIN "product"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "owner_id", "product_id", "title", "price",
						"city_id", "count", "available_count", "delivery", "safe_deal", "saler_id"}).
						AddRow(uint64(1), uint64(1), uint64(1), "Car", uint64(111),
							uint64(1), uint32(1), uint32(1), true, true, uint64(1)))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			basketStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(basketStorage)

			_, err = basketStorage.GetOrdersInBasketByUserID(ctx, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetOrdersNotInBasket(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city_id, "order".count, "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM public."order"
    INNER JOIN "product"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "owner_id", "product_id", "title", "price",
						"city_id", "count", "available_count", "delivery", "safe_deal", "saler_id"}).
						AddRow(uint64(1), uint64(1), uint64(1), "Car", uint64(111),
							uint64(1), uint32(1), uint32(1), true, true, uint64(1)))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			basketStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(basketStorage)

			_, err = basketStorage.GetOrdersNotInBasketByUserID(ctx, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetOrdersSold(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT "order".id, "order".owner_id, "order".product_id,
        "product".title, "product".price, "product".city_id, "order".count, "product".available_count,
        "product".delivery, "product".safe_deal, "product".saler_id FROM public."order"
    INNER JOIN "product"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "owner_id", "product_id", "title", "price",
						"city_id", "count", "available_count", "delivery", "safe_deal", "saler_id"}).
						AddRow(uint64(1), uint64(1), uint64(1), "Car", uint64(111),
							uint64(1), uint32(1), uint32(1), true, true, uint64(1)))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			basketStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(basketStorage)

			_, err = basketStorage.GetOrdersSoldByUserID(ctx, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

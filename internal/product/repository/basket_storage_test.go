package repository_test

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/pashagolub/pgxmock/v3"
	"testing"
	"time"
)

func TestDeleteOrder(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		ownerID                uint64
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
			ownerID: 1,
			orderID: 1,
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

			err = catStorage.DeleteOrder(ctx, testCase.orderID, testCase.ownerID)
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

	_ = my_logger.NewNop()

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

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		ownerID                uint64
		orderID                uint64
		newStatus              uint8
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
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
			behaviorProductStorage: func(m *repository.ProductStorage) {
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
			ctx := context.Background()

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage)

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

	_ = my_logger.NewNop()

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
       delivery, safe_deal, is_active FROM public."product"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"saler_id", "category_id", "title", "description", "price",
						"created_at", "views", "available_count", "city_id", "delivery", "safe_deal", "is_active"}).
						AddRow(uint64(1), uint64(1), "Car", "text", uint64(1212), time.Now(),
							uint32(6), uint32(4), uint64(6), true, true, true))

				mockPool.ExpectQuery(`SELECT url FROM public."image"`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"url"}).
						AddRow("safsafddasf"))

				mockPool.ExpectQuery(`SELECT COUNT`).WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{`count`}).
						AddRow(uint64(1)))

				mockPool.ExpectQuery(`SELECT id FROM public.favourite`).WithArgs(uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

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

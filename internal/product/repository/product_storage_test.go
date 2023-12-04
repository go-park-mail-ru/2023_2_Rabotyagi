package repository_test

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/pashagolub/pgxmock/v3"
	"testing"
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

func TestActivateProduct(t *testing.T) {
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

func TestDeleteProduct(t *testing.T) {
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

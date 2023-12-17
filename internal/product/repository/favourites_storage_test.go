package repository_test

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/pashagolub/pgxmock/v3"
	"testing"
)

func TestAddToFavourites(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.ProductStorage)
		ownerID                 uint64
		productID               uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`INSERT INTO public."favourite"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:   1,
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

			testCase.behaviorCategoryStorage(catStorage)

			err = catStorage.AddToFavourites(ctx, testCase.ownerID, testCase.productID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteFromFavourites(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.ProductStorage)
		ownerID                 uint64
		productID               uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."favourite"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:   1,
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

			testCase.behaviorCategoryStorage(catStorage)

			err = catStorage.DeleteFromFavourites(ctx, testCase.ownerID, testCase.productID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

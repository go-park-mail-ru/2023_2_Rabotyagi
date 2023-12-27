package repository_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/pashagolub/pgxmock/v3"
)

func TestDeleteFromFavourites(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		ownerID                 uint64
		productID               uint64
		expectedError           error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."favourite"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:       1,
			productID:     1,
			expectedError: nil,
		},
		{
			name: "test no affected rows",
			behaviorCategoryStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."favourite"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			ownerID:       1,
			productID:     1,
			expectedError: repository.ErrNoAffectedFavouriteRows,
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

			testCase.behaviorCategoryStorage(catStorage, mockPool)

			errActual := catStorage.DeleteFromFavourites(ctx, testCase.ownerID, testCase.productID)

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualError(errActual, testCase.expectedError)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAddToFavourites(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		ownerID                 uint64
		productID               uint64
		expectedError           error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`INSERT INTO public."favourite"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			ownerID:       1,
			productID:     1,
			expectedError: nil,
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

			testCase.behaviorCategoryStorage(catStorage, mockPool)

			errActual := catStorage.AddToFavourites(ctx, testCase.ownerID, testCase.productID)

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualError(errActual, testCase.expectedError)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

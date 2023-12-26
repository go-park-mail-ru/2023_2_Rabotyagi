package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

func TestAddPremium(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()
	beginPremium := time.Now()

	testError := myerrors.NewErrorInternal("test error")

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		userID                 uint64
		productID              uint64
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectExec(`UPDATE public."product" 
	SET premium_status=\$1, premium_begin=\$2,
    premium_expire=\$3 WHERE id=\$4 AND saler_id=\$5`).WithArgs(
					statuses.IntStatusPremiumSucceeded, beginPremium,
					beginPremium.AddDate(0, 0, 7), uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:        1,
			productID:     1,
			expectedError: nil,
		},
		{
			name: "test rowsAffected = 0",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectExec(`UPDATE public."product" 
	SET premium_status=\$1, premium_begin=\$2,
    premium_expire=\$3 WHERE id=\$4 AND saler_id=\$5`).WithArgs(
					statuses.IntStatusPremiumSucceeded, beginPremium,
					beginPremium.AddDate(0, 0, 7), uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			userID:        1,
			productID:     1,
			expectedError: repository.ErrNoAffectedProductRows,
		},
		{
			name: "test internal error",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectExec(`UPDATE public."product" 
	SET premium_status=\$1, premium_begin=\$2,
    premium_expire=\$3 WHERE id=\$4 AND saler_id=\$5`).WithArgs(
					statuses.IntStatusPremiumSucceeded, beginPremium,
					beginPremium.AddDate(0, 0, 7), uint64(1), uint64(1)).
					WillReturnError(testError)
				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			userID:        1,
			productID:     1,
			expectedError: testError,
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

			err = catStorage.AddPremium(ctx, testCase.productID, testCase.userID,
				beginPremium, beginPremium.AddDate(0, 0, 7))
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCheckPremium(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	testError := myerrors.NewErrorInternal("test error")

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		userID                 uint64
		productID              uint64
		expectedStatus         uint8
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectQuery(`SELECT premium_status FROM public."product"
                      WHERE id=\$1 AND saler_id=\$2`).WithArgs(uint64(1), uint64(1)).WillReturnRows(
					pgxmock.NewRows([]string{"premium_status"}).AddRow(uint8(1)))
			},
			userID:         1,
			productID:      1,
			expectedStatus: 1,
			expectedError:  nil,
		},
		{
			name: "test pgx.ErrNoRows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectQuery(`SELECT premium_status FROM public."product"
                      WHERE id=\$1 AND saler_id=\$2`).WithArgs(uint64(1), uint64(1)).WillReturnError(pgx.ErrNoRows)
			},
			userID:         1,
			productID:      1,
			expectedStatus: 0,
			expectedError:  repository.ErrPremiumStatusNotFound,
		},
		{
			name: "test internal error",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectQuery(`SELECT premium_status FROM public."product"
                      WHERE id=\$1 AND saler_id=\$2`).WithArgs(uint64(1), uint64(1)).WillReturnError(testError)
			},
			userID:         1,
			productID:      1,
			expectedStatus: 0,
			expectedError:  testError,
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

			receivedStatus, err := catStorage.CheckPremiumStatus(ctx, testCase.productID, testCase.userID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(receivedStatus, testCase.expectedStatus); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

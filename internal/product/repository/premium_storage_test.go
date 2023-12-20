package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/pashagolub/pgxmock/v3"
)

func TestAddPremium(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()
	beginPremium := time.Now()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
		productID              uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()
				mockPool.ExpectExec(`UPDATE public."product"`).WithArgs(beginPremium,
					beginPremium.AddDate(0, 0, 7), uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:    1,
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

			err = catStorage.AddPremium(ctx, testCase.productID, testCase.userID,
				beginPremium, beginPremium.AddDate(0, 0, 7))
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

package repository_test

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/pashagolub/pgxmock/v3"
	"testing"
	"time"
)

func TestGetUserWithoutPasswordByID(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                string
		behaviorCityStorage func(m *repository.UserStorage)
		userID              uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCityStorage: func(m *repository.UserStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT email, phone, name, birthday, avatar, created_at FROM public."user"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"email", "phone", "name", "birthday", "avatar", "created_at"}).
						AddRow("test@gmail.com", sql.NullString{String: "88005553535", Valid: true},
							sql.NullString{String: "John", Valid: true},
							testTime, sql.NullString{String: "afsghga", Valid: true}, testTime))

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

			cityStorage, err := repository.NewUserStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCityStorage(cityStorage)

			_, err = cityStorage.GetUserWithoutPasswordByID(ctx, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                string
		behaviorCityStorage func(m *repository.UserStorage)
		userID              uint64
		updDataMap          map[string]interface{}
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCityStorage: func(m *repository.UserStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."user"`).WithArgs(
					"Alex", uint64(1)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectQuery(`SELECT email, phone, name, birthday, avatar, created_at FROM public."user"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"email", "phone", "name", "birthday", "avatar", "created_at"}).
						AddRow("test@gmail.com", sql.NullString{String: "88005553535", Valid: true},
							sql.NullString{String: "John", Valid: true},
							testTime, sql.NullString{String: "afsghga", Valid: true}, testTime))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:     1,
			updDataMap: map[string]interface{}{"name": "Alex"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			cityStorage, err := repository.NewUserStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCityStorage(cityStorage)

			_, err = cityStorage.UpdateUser(ctx, testCase.userID, testCase.updDataMap)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

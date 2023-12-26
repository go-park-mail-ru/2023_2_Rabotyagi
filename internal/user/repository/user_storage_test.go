package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/pashagolub/pgxmock/v3"
)

func TestGetUserWithoutPasswordByID(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                string
		behaviorCityStorage func(m *repository.UserStorage, mockPool pgxmock.PgxPoolIface)
		userID              uint64
		expectedResponse    *models.UserWithoutPassword
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCityStorage: func(m *repository.UserStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT email, phone, name, birthday, avatar, created_at FROM public."user"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"email", "phone", "name", "birthday", "avatar", "created_at"}).
						AddRow("test@gmail.com", sql.NullString{String: "88005553535", Valid: true},
							sql.NullString{String: "John", Valid: true},
							testTime, sql.NullString{String: "afsghga", Valid: true}, testTime))

				mockPool.ExpectQuery(`SELECT AVG`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"avg"}).
						AddRow(sql.NullFloat64{Valid: false, Float64: 0}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID: 1,
			expectedResponse: &models.UserWithoutPassword{ //nolint:exhaustruct
				ID:        1,
				Email:     "test@gmail.com",
				Phone:     sql.NullString{String: "88005553535", Valid: true},
				Name:      sql.NullString{String: "John", Valid: true},
				Birthday:  sql.NullTime{Valid: true, Time: testTime},
				Avatar:    sql.NullString{String: "afsghga", Valid: true},
				CreatedAt: testTime,
			},
		},
		{
			name: "test not null rating",
			behaviorCityStorage: func(m *repository.UserStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT email, phone, name, birthday, avatar, created_at FROM public."user"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"email", "phone", "name", "birthday", "avatar", "created_at"}).
						AddRow("test@gmail.com", sql.NullString{String: "88005553535", Valid: true},
							sql.NullString{String: "John", Valid: true},
							testTime, sql.NullString{String: "afsghga", Valid: true}, testTime))

				mockPool.ExpectQuery(`SELECT AVG`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"avg"}).
						AddRow(sql.NullFloat64{Valid: true, Float64: 4}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID: 1,
			expectedResponse: &models.UserWithoutPassword{
				ID:        1,
				Email:     "test@gmail.com",
				AvgRating: sql.NullFloat64{Valid: true, Float64: 4},
				Phone:     sql.NullString{String: "88005553535", Valid: true},
				Name:      sql.NullString{String: "John", Valid: true},
				Birthday:  sql.NullTime{Valid: true, Time: testTime},
				Avatar:    sql.NullString{String: "afsghga", Valid: true},
				CreatedAt: testTime,
			},
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

			cityStorage, err := repository.NewUserStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCityStorage(cityStorage, mockPool)

			response, err := cityStorage.GetUserWithoutPasswordByID(ctx, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualTest(response, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                string
		behaviorCityStorage func(m *repository.UserStorage, mockPool pgxmock.PgxPoolIface)
		userID              uint64
		updDataMap          map[string]interface{}
		expectedResponse    *models.UserWithoutPassword
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCityStorage: func(m *repository.UserStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."user"`).WithArgs(
					"Alex", uint64(1)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mockPool.ExpectQuery(`SELECT email, phone, name, birthday, avatar, created_at FROM public."user"`).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"email", "phone", "name", "birthday", "avatar", "created_at"}).
						AddRow("test@gmail.com", sql.NullString{String: "88005553535", Valid: true},
							sql.NullString{String: "Alex", Valid: true},
							testTime, sql.NullString{String: "afsghga", Valid: true}, testTime))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:     1,
			updDataMap: map[string]interface{}{"name": "Alex"},
			expectedResponse: &models.UserWithoutPassword{ //nolint:exhaustruct
				ID:        1,
				Email:     "test@gmail.com",
				Phone:     sql.NullString{String: "88005553535", Valid: true},
				Name:      sql.NullString{String: "Alex", Valid: true},
				Birthday:  sql.NullTime{Valid: true, Time: testTime},
				Avatar:    sql.NullString{String: "afsghga", Valid: true},
				CreatedAt: testTime,
			},
		},
		{
			name: "test err no rows",
			behaviorCityStorage: func(m *repository.UserStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."user"`).WithArgs(
					"updated@mail.ru", uint64(1)).WillReturnError(sql.ErrNoRows)

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			userID:           1,
			updDataMap:       map[string]interface{}{"email": "updated@mail.ru"},
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

			userStorage, err := repository.NewUserStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCityStorage(userStorage, mockPool)

			response, _ := userStorage.UpdateUser(ctx, testCase.userID, testCase.updDataMap)

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualTest(response, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

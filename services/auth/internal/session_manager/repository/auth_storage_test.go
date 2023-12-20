package repository_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/repository"

	"github.com/pashagolub/pgxmock/v3"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                string
		behaviorAuthStorage func(m *repository.AuthStorage)
		testEmail           string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorAuthStorage: func(m *repository.AuthStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT id FROM public."user"`).WithArgs("test@gmail.com").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).
						AddRow("1"))

				mockPool.ExpectQuery(`SELECT id, email, password FROM public."user"`).WithArgs("test@gmail.com").
					WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password"}).
						AddRow(uint64(1), "test@gmail.com", "123456"))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			testEmail: "test@gmail.com",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewAuthStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorAuthStorage(catStorage)

			_, err = catStorage.GetUser(ctx, testCase.testEmail)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                string
		behaviorAuthStorage func(m *repository.AuthStorage)
		testEmail           string
		testPassword        string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorAuthStorage: func(m *repository.AuthStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT id FROM public."user"`).WithArgs("test@gmail.com").
					WillReturnRows(pgxmock.NewRows([]string{"id"}))

				mockPool.ExpectExec(`INSERT INTO public."user"`).WithArgs("test@gmail.com", "123456").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectQuery(`SELECT last_value FROM "public"."user_id_seq";`).
					WillReturnRows(pgxmock.NewRows([]string{"last_value"}).
						AddRow(uint64(1)))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			testEmail:    "test@gmail.com",
			testPassword: "123456",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewAuthStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorAuthStorage(catStorage)

			_, err = catStorage.AddUser(ctx, testCase.testEmail, testCase.testPassword)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

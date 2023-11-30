package repository_test

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/jackc/pgx/v5"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func TestGetFullCategories(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.CategoryStorage)
		expectedResponse        any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.CategoryStorage) {
				mockPool.ExpectBeginTx(pgx.TxOptions{})
				mockPool.ExpectQuery(`^SELECT (.+) FROM public."category"$`).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "parent_id"}).
						AddRow(1, "Animal", 4).
						AddRow(2, "Cats", 1).
						AddRow(3, "Dogs", 1))
				mockPool.ExpectCommit()
			},
			expectedResponse: []*models.Category{
				{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 4}},
				{ID: 2, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
				{ID: 3, Name: "Dogs", ParentID: sql.NullInt64{Valid: true, Int64: 1}}},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			catStorage, err := repository.NewCategoryStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCategoryStorage(catStorage)

			response, err := catStorage.GetFullCategories(ctx)
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

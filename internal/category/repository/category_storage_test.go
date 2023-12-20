package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/pashagolub/pgxmock/v3"
)

func TestGetFullCategories(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

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
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`^SELECT (.+) FROM public."category"$`).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "parent_id"}).
						AddRow(uint64(1), "Animal", sql.NullInt64{Valid: false, Int64: 0}).
						AddRow(uint64(2), "Cats", sql.NullInt64{Valid: true, Int64: 1}).
						AddRow(uint64(3), "Dogs", sql.NullInt64{Valid: true, Int64: 1}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.Category{
				{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
				{ID: 2, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
				{ID: 3, Name: "Dogs", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

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

func TestSearchCategory(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.CategoryStorage)
		expectedResponse        any
		searchInput             string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.CategoryStorage) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT category.id, category.name, category.parent_id FROM public."category"`).
					WithArgs("%ani%").
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "parent_id"}).
						AddRow(uint64(1), "Animal", sql.NullInt64{Valid: false, Int64: 0}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.Category{
				{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
			},
			searchInput: "Ani",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewCategoryStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCategoryStorage(catStorage)

			response, err := catStorage.SearchCategory(ctx, testCase.searchInput)
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

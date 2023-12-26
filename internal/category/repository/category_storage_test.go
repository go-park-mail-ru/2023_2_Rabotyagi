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

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface)
		expectedResponse        any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`^SELECT (.+) FROM public."category"$`).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "parent_id"}).
						AddRow(uint64(1), "Animal", sql.NullInt64{Valid: false, Int64: 0}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.Category{
				{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
			},
		},
		{
			name: "test empty",
			behaviorCategoryStorage: func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`^SELECT (.+) FROM public."category"$`).
					WillReturnRows(pgxmock.NewRows([]string{}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.Category(nil),
		},
		{
			name: "test more data",
			behaviorCategoryStorage: func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface) {
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

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewCategoryStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCategoryStorage(catStorage, mockPool)

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

	type TestCase struct {
		name                    string
		behaviorCategoryStorage func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface)
		expectedResponse        any
		searchInput             string
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface) {
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
		{
			name: "test empty",
			behaviorCategoryStorage: func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT category.id, category.name, category.parent_id FROM public."category"`).
					WithArgs("%ani%").
					WillReturnRows(pgxmock.NewRows([]string{}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.Category(nil),
			searchInput:      "Ani",
		},
		{
			name: "test more data",
			behaviorCategoryStorage: func(m *repository.CategoryStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()
				mockPool.ExpectQuery(`SELECT category.id, category.name, category.parent_id FROM public."category"`).
					WithArgs("%s%").
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "parent_id"}).
						AddRow(uint64(1), "Shoes", sql.NullInt64{Valid: false, Int64: 0}).
						AddRow(uint64(2), "Sunglasses", sql.NullInt64{Valid: true, Int64: 1}).
						AddRow(uint64(3), "Sportswear", sql.NullInt64{Valid: true, Int64: 1}))
				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			expectedResponse: []*models.Category{
				{ID: 1, Name: "Shoes", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
				{ID: 2, Name: "Sunglasses", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
				{ID: 3, Name: "Sportswear", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
			},
			searchInput: "S",
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

			catStorage, err := repository.NewCategoryStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorCategoryStorage(catStorage, mockPool)

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

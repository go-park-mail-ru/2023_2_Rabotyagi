package repository_test

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
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
		behaviorCategoryStorage func(m *mocks.MockICategoryStorage)
		expectedResponse        any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCategoryStorage: func(m *mocks.MockICategoryStorage) {
				m.EXPECT().GetFullCategories(gomock.Any()).Return([]*models.Category{
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 2, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
					{ID: 3, Name: "Dogs", ParentID: sql.NullInt64{Valid: true, Int64: 1}}}, nil)

				mockPool.ExpectQuery(`SELECT "category".id,"category".name, "category".parent_id FROM public."category"`).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "parent_id"}).AddRow(1, "Animal", nil).
						AddRow(2, "Cats", 1).
						AddRow(3, "Dogs", 1))

			},
			expectedResponse: []*models.Category{
				{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
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

			catStorage := mocks.NewMockICategoryStorage(ctrl)

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
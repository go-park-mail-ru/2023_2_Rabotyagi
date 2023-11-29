package delivery

import (
	"database/sql"
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

//nolint:funlen
func TestGetFullCategories(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		method           string
		behavior         func(m *mocks.MockICategoryService)
		expectedResponse CategoryListResponse
	}

	testCases := [...]TestCase{
		{
			name:   "test basic work",
			method: http.MethodGet,
			behavior: func(m *mocks.MockICategoryService) {
				m.EXPECT().GetFullCategories(gomock.Any()).Return([]*models.Category{
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
					{ID: 3, Name: "Dogs", ParentID: sql.NullInt64{Valid: true, Int64: 1}}}, nil)
			},
			expectedResponse: CategoryListResponse{Status: statuses.StatusResponseSuccessful,
				Body: []*models.Category{{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 1}},
					{ID: 3, Name: "Dogs", ParentID: sql.NullInt64{Valid: true, Int64: 1}}}},
		},
		{
			name: "test empty work",
			behavior: func(m *mocks.MockICategoryService) {
				m.EXPECT().GetFullCategories(gomock.Any()).Return([]*models.Category{}, nil)
			},
			expectedResponse: CategoryListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []*models.Category{},
			},
		},
		{
			name: "test repeated names a lot",
			behavior: func(m *mocks.MockICategoryService) {
				m.EXPECT().GetFullCategories(gomock.Any()).Return(
					[]*models.Category{
						{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
						{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
						{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
						{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
						{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
						{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					}, nil)
			},
			expectedResponse: CategoryListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: []*models.Category{
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 1, Name: "Animal", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServ := mocks.NewMockICategoryService(ctrl)

			testCase.behavior(mockServ)

			cityHandler, err := NewCategoryHandler(mockServ)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			req := httptest.NewRequest(testCase.method, "/api/v1/city/get_full", nil)

			w := httptest.NewRecorder()

			cityHandler.GetFullCategories(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse CategoryListResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			if !reflect.DeepEqual(testCase.expectedResponse, resultResponse) {
				t.Errorf("Wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}

func TestSearchCategoryHandler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		searchInput      string
		method           string
		behavior         func(m *mocks.MockICategoryService)
		expectedResponse CategoryListResponse
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			method:      http.MethodGet,
			searchInput: "Ca",
			behavior: func(m *mocks.MockICategoryService) {
				m.EXPECT().SearchCategory(gomock.Any(), "Ca").Return([]*models.Category{
					{ID: 3, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 2}},
					{ID: 7, Name: "Cars", ParentID: sql.NullInt64{Valid: true, Int64: 4}}}, nil)
			},
			expectedResponse: CategoryListResponse{Status: statuses.StatusResponseSuccessful,
				Body: []*models.Category{{ID: 3, Name: "Cats", ParentID: sql.NullInt64{Valid: true, Int64: 2}},
					{ID: 7, Name: "Cars", ParentID: sql.NullInt64{Valid: true, Int64: 4}}}},
		},
		{
			name:        "test empty query",
			method:      http.MethodGet,
			searchInput: "",
			behavior: func(m *mocks.MockICategoryService) {
				m.EXPECT().SearchCategory(gomock.Any(), "").Return([]*models.Category{}, nil)
			},
			expectedResponse: CategoryListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []*models.Category{},
			},
		},
		{
			name:        "test special symbols query",
			searchInput: "Кошки &&& Собаки",
			behavior: func(m *mocks.MockICategoryService) {
				m.EXPECT().SearchCategory(gomock.Any(), "Кошки &&& Собаки").Return([]*models.Category{
					{ID: 1, Name: "Кошки", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 2, Name: "Собаки", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
				}, nil)
			},
			expectedResponse: CategoryListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: []*models.Category{
					{ID: 1, Name: "Кошки", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
					{ID: 2, Name: "Собаки", ParentID: sql.NullInt64{Valid: false, Int64: 0}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServ := mocks.NewMockICategoryService(ctrl)

			testCase.behavior(mockServ)

			cityHandler, err := NewCategoryHandler(mockServ)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			req := httptest.NewRequest(testCase.method, `/api/v1/city/search`, nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"searched": testCase.searchInput})

			w := httptest.NewRecorder()

			cityHandler.SearchCategoryHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse CategoryListResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			if !reflect.DeepEqual(testCase.expectedResponse, resultResponse) {
				t.Errorf("Wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}

package delivery_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"

	"go.uber.org/mock/gomock"
)

func TestGetFullCities(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		behavior         func(m *mocks.MockICityService)
		expectedResponse any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behavior: func(m *mocks.MockICityService) {
				m.EXPECT().GetFullCities(gomock.Any()).Return([]*models.City{
					{ID: 1, Name: "Moscow"}, {ID: 2, Name: "Voronezh"},
				}, nil)
			},
			expectedResponse: delivery.CityListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: []*models.City{
					{ID: 1, Name: "Moscow"}, {ID: 2, Name: "Voronezh"},
				},
			},
		},
		{
			name: "test empty work",
			behavior: func(m *mocks.MockICityService) {
				m.EXPECT().GetFullCities(gomock.Any()).Return([]*models.City{}, nil)
			},
			expectedResponse: delivery.CityListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []*models.City{},
			},
		},
		{
			name: "test repeated names a lot",
			behavior: func(m *mocks.MockICityService) {
				m.EXPECT().GetFullCities(gomock.Any()).Return(
					[]*models.City{
						{ID: 1, Name: "Moscow"},
						{ID: 1, Name: "Moscow"},
						{ID: 1, Name: "Moscow"},
						{ID: 1, Name: "Moscow"},
						{ID: 1, Name: "Moscow"},
						{ID: 1, Name: "Moscow"},
					}, nil)
			},
			expectedResponse: delivery.CityListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: []*models.City{
					{ID: 1, Name: "Moscow"},
					{ID: 1, Name: "Moscow"},
					{ID: 1, Name: "Moscow"},
					{ID: 1, Name: "Moscow"},
					{ID: 1, Name: "Moscow"},
					{ID: 1, Name: "Moscow"},
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

			mockServ := mocks.NewMockICityService(ctrl)

			testCase.behavior(mockServ)

			cityHandler, err := delivery.NewCityHandler(mockServ)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/city/get_full", nil)

			w := httptest.NewRecorder()

			cityHandler.GetFullCitiesHandler(w, req)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestSearchCityHandler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		query            string
		behavior         func(m *mocks.MockICityService)
		expectedResponse any
	}

	testCases := [...]TestCase{
		{
			name:  "test basic work",
			query: "Москва",
			behavior: func(m *mocks.MockICityService) {
				m.EXPECT().SearchCity(gomock.Any(), "Москва").Return([]*models.City{
					{ID: 1, Name: "Moscow"},
				}, nil)
			},
			expectedResponse: delivery.CityListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: []*models.City{
					{ID: 1, Name: "Moscow"},
				},
			},
		},
		{
			name:  "test empty query",
			query: "",
			behavior: func(m *mocks.MockICityService) {
				m.EXPECT().SearchCity(gomock.Any(), "").Return([]*models.City{}, nil)
			},
			expectedResponse: delivery.CityListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []*models.City{},
			},
		},
		{
			name:  "test special symbols query",
			query: "Москва &&& Вологда",
			behavior: func(m *mocks.MockICityService) {
				m.EXPECT().SearchCity(gomock.Any(), "Москва &&& Вологда").Return([]*models.City{
					{ID: 1, Name: "Вологда"}, {ID: 1, Name: "Вологда"},
				}, nil)
			},
			expectedResponse: delivery.CityListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: []*models.City{
					{ID: 1, Name: "Вологда"}, {ID: 1, Name: "Вологда"},
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

			mockServ := mocks.NewMockICityService(ctrl)

			testCase.behavior(mockServ)

			cityHandler, err := delivery.NewCityHandler(mockServ)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/city/search", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"searched": testCase.query})
			w := httptest.NewRecorder()

			cityHandler.SearchCityHandler(w, req)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

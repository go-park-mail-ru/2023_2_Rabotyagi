package delivery_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"go.uber.org/mock/gomock"
)

//nolint:funlen
func TestGetFullCities(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		behavior         func(m *mocks.MockICityService)
		expectedResponse delivery.CityListResponse
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

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse delivery.CityListResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			err = utils.EqualTest(resultResponse, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

//nolint:funlen
func TestSearchCityHandler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		query            string
		behavior         func(m *mocks.MockICityService)
		expectedResponse delivery.CityListResponse
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

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse delivery.CityListResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			err = utils.EqualTest(resultResponse, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

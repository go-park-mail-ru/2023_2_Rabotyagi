package delivery_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/city/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"

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
					{ID: 1, Name: "Moscow"}, {ID: 2, Name: "Voronezh"}}, nil)
			},
			expectedResponse: delivery.CityListResponse{Status: statuses.StatusResponseSuccessful,
				Body: []*models.City{{ID: 1, Name: "Moscow"}, {ID: 2, Name: "Voronezh"}}},
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

			if !reflect.DeepEqual(testCase.expectedResponse, resultResponse) {
				t.Errorf("Wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}

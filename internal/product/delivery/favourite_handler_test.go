package delivery_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	mocksauth "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddToFavourites(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                     string
		behaviorFavouriteService func(m *mocks.MockIProductService)
		request                  *http.Request
		expectedResponse         responses.ResponseID
	}
	testCases := [...]TestCase{
		{
			name:    "test basic work",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add-to-fav", strings.NewReader(`{"product_id":1}`)),
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddToFavourites(gomock.Any(), uint64(testUserID), io.NopCloser(strings.NewReader(
					`{"product_id":1}`))).Return(nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProductService := mocks.NewMockIProductService(ctrl)
			mockSessionManagerClient := mocksauth.NewMockSessionMangerClient(ctrl)

			behaviorSessionManagerClientCheck(mockSessionManagerClient)
			testCase.behaviorFavouriteService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&testCookie)
			productHandler.AddToFavouritesHandler(w, testCase.request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.ResponseID

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

func TestGetFavourites(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                     string
		behaviorFavouriteService func(m *mocks.MockIProductService)
		request                  *http.Request
		expectedResponse         *delivery.ProductListResponse
	}
	testCases := [...]TestCase{
		{
			name:    "test basic work",
			request: httptest.NewRequest(http.MethodGet, "/api/v1/profile/favourites", nil),
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetUserFavourites(gomock.Any(), uint64(testUserID)).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "sofa"}, {ID: 2, Title: "laptop"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse([]*models.ProductInFeed{{ID: 1, Title: "sofa"}, {ID: 2, Title: "laptop"}}),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProductService := mocks.NewMockIProductService(ctrl)
			mockSessionManagerClient := mocksauth.NewMockSessionMangerClient(ctrl)

			behaviorSessionManagerClientCheck(mockSessionManagerClient)
			testCase.behaviorFavouriteService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&testCookie)
			productHandler.GetFavouritesHandler(w, testCase.request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse delivery.ProductListResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			err = utils.EqualTest(&resultResponse, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteFavourite(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                     string
		behaviorFavouriteService func(m *mocks.MockIProductService)
		expectedResponse         responses.ResponseID
		queryProductID           string
	}
	testCases := [...]TestCase{
		{
			name:           "test basic work",
			queryProductID: "1",
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().DeleteFromFavourites(gomock.Any(), uint64(testUserID), uint64(1)).Return(nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProductService := mocks.NewMockIProductService(ctrl)
			mockSessionManagerClient := mocksauth.NewMockSessionMangerClient(ctrl)

			behaviorSessionManagerClientCheck(mockSessionManagerClient)
			testCase.behaviorFavouriteService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodDelete, "/api/v1/product/remove-from-fav", nil)
			utils.AddQueryParamsToRequest(request, map[string]string{"product_id": testCase.queryProductID})

			request.AddCookie(&testCookie)
			productHandler.DeleteFromFavouritesHandler(w, request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.ResponseID

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

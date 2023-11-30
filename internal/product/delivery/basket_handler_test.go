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

func TestAddOrder(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       *delivery.OrderResponse
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/order/add", strings.NewReader(
				`{
					"owner_id": 67890,
					"saler_id": 54321,
					"product_id": 98765,
					"city_id": 24680,
					"title": "Example Product",
					"price": 999,
					"count": 3,
					"available_count": 5,
					"delivery": true,
					"safe_deal": false,
					"in_favourites": true,
					"images": [{"url":"img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695"},
								{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddOrder(gomock.Any(), io.NopCloser(strings.NewReader(
					`{
					"owner_id": 67890,
					"saler_id": 54321,
					"product_id": 98765,
					"city_id": 24680,
					"title": "Example Product",
					"price": 999,
					"count": 3,
					"available_count": 5,
					"delivery": true,
					"safe_deal": false,
					"in_favourites": true,
					"images": [{"url":"img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695"},
								{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)),
					uint64(testUserID)).Return(&models.OrderInBasket{
					OwnerID:        67890,
					SalerID:        54321,
					ProductID:      98765,
					CityID:         24680,
					Title:          "Example Product",
					Price:          999,
					Count:          3,
					AvailableCount: 5,
					Delivery:       true,
					SafeDeal:       false,
					InFavourites:   true,
					Images: []models.Image{
						{
							URL: "img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695",
						},
						{
							URL: "8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1",
						},
					},
				}, nil)
			},
			expectedResponse: &delivery.OrderResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: &models.OrderInBasket{
					OwnerID:        67890,
					SalerID:        54321,
					ProductID:      98765,
					CityID:         24680,
					Title:          "Example Product",
					Price:          999,
					Count:          3,
					AvailableCount: 5,
					Delivery:       true,
					SafeDeal:       false,
					InFavourites:   true,
					Images: []models.Image{
						{
							URL: "img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695",
						},
						{
							URL: "8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1",
						},
					},
				},
			},
		},
		{
			name:    "test partial",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/order/add", strings.NewReader(`{"product_id":3, "count":3}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddOrder(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"product_id":3, "count":3}`)), uint64(testUserID)).Return(&models.OrderInBasket{
					ProductID: 3,
					Count:     3,
				}, nil)
			},
			expectedResponse: &delivery.OrderResponse{
				Status: statuses.StatusResponseSuccessful,
				Body: &models.OrderInBasket{
					ProductID: 3,
					Count:     3,
				},
			},
		},
		{
			name:    "test empty",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/order/add", strings.NewReader(`{}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddOrder(gomock.Any(), io.NopCloser(strings.NewReader(
					`{}`)), uint64(testUserID)).Return(&models.OrderInBasket{}, nil)
			},
			expectedResponse: &delivery.OrderResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   &models.OrderInBasket{},
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
			testCase.behaviorProductService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&testCookie)
			productHandler.AddOrderHandler(w, testCase.request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse *delivery.OrderResponse

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

func TestGetBasket(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                     string
		behaviorFavouriteService func(m *mocks.MockIProductService)
		request                  *http.Request
		expectedResponse         any
	}
	testCases := [...]TestCase{
		{
			name:    "test basic work",
			request: httptest.NewRequest(http.MethodGet, "/api/v1/order/get_basket", nil),
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetOrdersByUserID(gomock.Any(), uint64(testUserID)).Return(
					[]*models.OrderInBasket{{ProductID: 1, Title: "sofa"}, {ProductID: 2, Title: "laptop"}}, nil)
			},
			expectedResponse: delivery.NewOrderListResponse([]*models.OrderInBasket{{ProductID: 1, Title: "sofa"}, {ProductID: 2, Title: "laptop"}}),
		},
		{
			name:    "test basic work",
			request: httptest.NewRequest(http.MethodGet, "/api/v1/order/get_basket", nil),
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetOrdersByUserID(gomock.Any(), uint64(testUserID)).Return(
					[]*models.OrderInBasket{}, nil)
			},
			expectedResponse: delivery.NewOrderListResponse([]*models.OrderInBasket{}),
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
			productHandler.GetBasketHandler(w, testCase.request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			expectedResponseRaw, err := json.Marshal(testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed to json.Marshal testCase.expectedResponse: %v", err)
			}

			err = utils.EqualTest(receivedResponse, expectedResponseRaw)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateOrderCountBasket(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                     string
		behaviorFavouriteService func(m *mocks.MockIProductService)
		request                  *http.Request
		expectedResponse         any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/order/update_count",
				strings.NewReader(`{"id":3, "count":3}`)),
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateOrderCount(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"id":3, "count":3}`)), uint64(testUserID)).Return(nil)
			},
			expectedResponse: responses.NewResponseSuccessful(delivery.ResponseSuccessfulUpdateCountOrder),
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
			productHandler.UpdateOrderCountHandler(w, testCase.request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			expectedResponseRaw, err := json.Marshal(testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed to json.Marshal testCase.expectedResponse: %v", err)
			}

			err = utils.EqualTest(receivedResponse, expectedResponseRaw)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateOrderStatusBasket(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                     string
		behaviorFavouriteService func(m *mocks.MockIProductService)
		request                  *http.Request
		expectedResponse         any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/order/update_status",
				strings.NewReader(`{"id":3, "status":1}`)),
			behaviorFavouriteService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateOrderStatus(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"id":3, "status":1}`)), uint64(testUserID)).Return(nil)
			},
			expectedResponse: responses.NewResponseSuccessful(delivery.ResponseSuccessfulUpdateStatusOrder),
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
			productHandler.UpdateOrderStatusHandler(w, testCase.request)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			expectedResponseRaw, err := json.Marshal(testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed to json.Marshal testCase.expectedResponse: %v", err)
			}

			err = utils.EqualTest(receivedResponse, expectedResponseRaw)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
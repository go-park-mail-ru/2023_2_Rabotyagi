package delivery_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	mocksauth "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"

	"go.uber.org/mock/gomock"
)

// testAccessToken for read only, because async usage.
const testAccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJlbWFpbCI6Iml2bi0xNS0wN0BtYWlsLnJ1IiwiZXhwaXJlIjoxNzAxMjg1MzE4LCJ1c2VySUQiOjExfQ." +
	"jIPlwcF5xGPpgQ5WYp5kFv9Av-yguX2aOYsAgbodDM4"

// testCookie for read only, because async usage.
var testCookie = http.Cookie{
	Name:  responses.CookieAuthName,
	Value: testAccessToken, Expires: time.Now().Add(time.Hour),
}

const testUserID = 1

var behaviorSessionManagerClientCheck = func(m *mocksauth.MockSessionMangerClient) { //nolint:gochecknoglobals
	m.EXPECT().Check(gomock.Any(), &auth.Session{AccessToken: testAccessToken}).Return(
		&auth.UserID{UserId: testUserID}, nil)
}

//nolint:funlen
func TestAddProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       responses.ResponseID
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add", strings.NewReader(
				`{"saler_id":1,
"category_id" :2,
"title":"adsf",
"description":"  description",
"price":123,
"available_count":1,
"city_id":1,
"delivery":false, "safe_deal":false}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"saler_id":1,
"category_id" :2,
"title":"adsf",
"description":"  description",
"price":123,
"available_count":1,
"city_id":1,
"delivery":false, "safe_deal":false}`)), uint64(testUserID)).Return(uint64(1), nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test another product",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add", strings.NewReader(
				`{"saler_id":1,
"category_id" :1,
"title":"TItle",
"description":"  de scription",
"price":1232,
"available_count":12,
"city_id":1,
"delivery":false, "safe_deal":false}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"saler_id":1,
"category_id" :1,
"title":"TItle",
"description":"  de scription",
"price":1232,
"available_count":12,
"city_id":1,
"delivery":false, "safe_deal":false}`)), uint64(testUserID)).Return(uint64(1), nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test product with images",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add", strings.NewReader(
				`{"saler_id":1,
"category_id" :1,
"title":"TItle",
"description":"  de scription",
"price":1232,
"available_count":12,
"city_id":1,
"delivery":false, "safe_deal":false, 
"images": [{"url":"img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695"},
{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"saler_id":1,
"category_id" :1,
"title":"TItle",
"description":"  de scription",
"price":1232,
"available_count":12,
"city_id":1,
"delivery":false, "safe_deal":false, 
"images": [{"url":"img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695"},
{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)), uint64(testUserID)).Return(uint64(1), nil)
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
			testCase.behaviorProductService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&testCookie)
			productHandler.AddProductHandler(w, testCase.request)

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

//nolint:funlen
func TestGetProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		idProduct              string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       *delivery.ProductResponse
	}

	testCases := [...]TestCase{
		{
			name:      "test basic work",
			idProduct: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProduct(gomock.Any(), uint64(1), uint64(testUserID)).Return(
					&models.Product{ID: 1, Title: "Title"}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductResponse(&models.Product{ID: 1, Title: "Title"}), //nolint:exhaustruct
		},
		{
			name:      "test empty product",
			idProduct: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProduct(gomock.Any(), uint64(1), uint64(testUserID)).Return(
					&models.Product{}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductResponse(&models.Product{}), //nolint:exhaustruct
		},
		{
			name:      "test full required fields of product",
			idProduct: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProduct(gomock.Any(), uint64(1), uint64(testUserID)).Return(
					&models.Product{ //nolint:exhaustruct
						ID: 1, SalerID: 1, CategoryID: 1, CityID: 1,
						Title: "Title", Description: "desc", Price: 123,
						CreatedAt: time.Unix(0, 0), Views: 12, AvailableCount: 1, Favourites: 12,
					}, nil)
			},
			expectedResponse: delivery.NewProductResponse(&models.Product{ //nolint:exhaustruct
				ID: 1, SalerID: 1, CategoryID: 1, CityID: 1,
				Title: "Title", Description: "desc", Price: 123,
				CreatedAt: time.Unix(0, 0), Views: 12, AvailableCount: 1, Favourites: 12,
			}),
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

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"id": testCase.idProduct})

			req.AddCookie(&testCookie)
			productHandler.GetProductHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse delivery.ProductResponse

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

//nolint:funlen
func TestGetProductList(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       *delivery.ProductListResponse
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "last_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsList(gomock.Any(), uint64(1), uint64(2), uint64(testUserID)).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "last_id": "0"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsList(gomock.Any(), uint64(0), uint64(0), uint64(testUserID)).Return(
					[]*models.ProductInFeed{}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "10", "last_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsList(gomock.Any(), uint64(1), uint64(10), uint64(testUserID)).Return(
					[]*models.ProductInFeed{
						{ID: 1, Title: "Title"},
						{ID: 2, Title: "Title2"},
						{ID: 3, Title: "Title2"},
						{ID: 4, Title: "Title2"},
						{ID: 5, Title: "Title2"},
						{ID: 6, Title: "Title2"},
						{ID: 7, Title: "Title2"},
						{ID: 8, Title: "Title2"},
						{ID: 9, Title: "Title2"},
						{ID: 10, Title: "Title2"},
					}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{
					{ID: 1, Title: "Title"},
					{ID: 2, Title: "Title2"},
					{ID: 3, Title: "Title2"},
					{ID: 4, Title: "Title2"},
					{ID: 5, Title: "Title2"},
					{ID: 6, Title: "Title2"},
					{ID: 7, Title: "Title2"},
					{ID: 8, Title: "Title2"},
					{ID: 9, Title: "Title2"},
					{ID: 10, Title: "Title2"},
				}), //nolint:exhaustruct
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

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_list", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)

			req.AddCookie(&testCookie)
			productHandler.GetProductListHandler(w, req)

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

//nolint:funlen
func TestGetListProductOfSaler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       *delivery.ProductListResponse
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "last_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(2), uint64(testUserID), true).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "last_id": "0"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(0), uint64(0), uint64(testUserID), true).Return(
					[]*models.ProductInFeed{}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "10", "last_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(10), uint64(testUserID), true).Return(
					[]*models.ProductInFeed{
						{ID: 1, Title: "Title"},
						{ID: 2, Title: "Title2"},
						{ID: 3, Title: "Title2"},
						{ID: 4, Title: "Title2"},
						{ID: 5, Title: "Title2"},
						{ID: 6, Title: "Title2"},
						{ID: 7, Title: "Title2"},
						{ID: 8, Title: "Title2"},
						{ID: 9, Title: "Title2"},
						{ID: 10, Title: "Title2"},
					}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{
					{ID: 1, Title: "Title"},
					{ID: 2, Title: "Title2"},
					{ID: 3, Title: "Title2"},
					{ID: 4, Title: "Title2"},
					{ID: 5, Title: "Title2"},
					{ID: 6, Title: "Title2"},
					{ID: 7, Title: "Title2"},
					{ID: 8, Title: "Title2"},
					{ID: 9, Title: "Title2"},
					{ID: 10, Title: "Title2"},
				}), //nolint:exhaustruct
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

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_list_of_saler", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)

			req.AddCookie(&testCookie)
			productHandler.GetListProductOfSalerHandler(w, req)

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

//nolint:funlen
func TestGetListProductOfAnotherSaler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       *delivery.ProductListResponse
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "last_id": "1", "saler_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(2), uint64(testUserID), false).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "last_id": "0", "saler_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(0), uint64(0), uint64(testUserID), false).Return(
					[]*models.ProductInFeed{}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "10", "last_id": "1", "saler_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(10), uint64(testUserID), false).Return(
					[]*models.ProductInFeed{
						{ID: 1, Title: "Title"},
						{ID: 2, Title: "Title2"},
						{ID: 3, Title: "Title2"},
						{ID: 4, Title: "Title2"},
						{ID: 5, Title: "Title2"},
						{ID: 6, Title: "Title2"},
						{ID: 7, Title: "Title2"},
						{ID: 8, Title: "Title2"},
						{ID: 9, Title: "Title2"},
						{ID: 10, Title: "Title2"},
					}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{
					{ID: 1, Title: "Title"},
					{ID: 2, Title: "Title2"},
					{ID: 3, Title: "Title2"},
					{ID: 4, Title: "Title2"},
					{ID: 5, Title: "Title2"},
					{ID: 6, Title: "Title2"},
					{ID: 7, Title: "Title2"},
					{ID: 8, Title: "Title2"},
					{ID: 9, Title: "Title2"},
					{ID: 10, Title: "Title2"},
				}), //nolint:exhaustruct
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

			testCase.behaviorProductService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_list_of_another_saler", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)

			productHandler.GetListProductOfAnotherSalerHandler(w, req)

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

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       responses.ResponseID
	}

	testCases := [...]TestCase{
		{
			name: "test basic work patch",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/product/update?id=1", strings.NewReader(
				`{"available_count": 1,
  "category_id": 1,
  "delivery": true,
  "description": "description not empty"}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"available_count": 1,
  "category_id": 1,
  "delivery": true,
  "description": "description not empty"}`)), true, uint64(1), uint64(testUserID)).Return(nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test empty patch",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/product/update?id=1", strings.NewReader(
				``)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					``)), true, uint64(1), uint64(testUserID)).Return(nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test basic work put",
			request: httptest.NewRequest(http.MethodPut, "/api/v1/product/update?id=1", strings.NewReader(
				`{"available_count": 1,
  "category_id": 1,  "city_id": 1, "saler_id": 1,
  "title": "title", "price" : 123,
  "description": "description not empty"}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"available_count": 1,
  "category_id": 1,  "city_id": 1, "saler_id": 1,
  "title": "title", "price" : 123,
  "description": "description not empty"}`)), false, uint64(1), uint64(testUserID)).Return(nil)
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
			testCase.behaviorProductService(mockProductService)

			productHandler, err := delivery.NewProductHandler(mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&testCookie)
			productHandler.UpdateProductHandler(w, testCase.request)

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

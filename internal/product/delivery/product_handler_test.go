package delivery_test

import (
	"fmt"
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
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
)

var behaviorSessionManagerClientCheck = func(m *mocksauth.MockSessionMangerClient) { //nolint:gochecknoglobals
	m.EXPECT().Check(gomock.Any(), &auth.Session{AccessToken: test.AccessToken}).Return(
		&auth.UserID{UserId: test.UserID}, nil).AnyTimes()
}

func NewProductHandler(ctrl *gomock.Controller,
	behaviorProductService func(m *mocks.MockIProductService),
) (*delivery.ProductHandler, error) {
	mockProductService := mocks.NewMockIProductService(ctrl)
	mockSessionManagerClient := mocksauth.NewMockSessionMangerClient(ctrl)

	behaviorSessionManagerClientCheck(mockSessionManagerClient)
	behaviorProductService(mockProductService)

	productHandler, err := delivery.NewProductHandler("test",
		"test", "test", "test",
		mockProductService, mockSessionManagerClient)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return productHandler, nil
}

func TestAddProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       any
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
"delivery":false, "safe_deal":false}`)), test.UserID).Return(uint64(1), nil)
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
"delivery":false, "safe_deal":false}`)), test.UserID).Return(uint64(1), nil)
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
{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)), test.UserID).Return(uint64(1), nil)
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
"delivery":false, "safe_deal":false}`)), test.UserID).Return(uint64(1), nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test error AddProduct",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add", strings.NewReader(
				``)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddProduct(gomock.Any(), io.NopCloser(strings.NewReader(
					``)), test.UserID).Return(uint64(0), myerrors.NewErrorInternal("Test internal error"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name: "test wrong method",
			request: httptest.NewRequest(http.MethodDelete, "/api/v1/product/add", strings.NewReader(
				``)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
			},
			expectedResponse: `Method not allowed
`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&test.Cookie)
			productHandler.AddProductHandler(w, testCase.request)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		idProduct              string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:      "test basic work",
			idProduct: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProduct(gomock.Any(), uint64(1), test.UserID).Return(
					&models.Product{ID: 1, Title: "Title"}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductResponse(&models.Product{ID: 1, Title: "Title"}), //nolint:exhaustruct
		},
		{
			name:      "test empty product",
			idProduct: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProduct(gomock.Any(), uint64(1), test.UserID).Return(
					&models.Product{}, nil) //nolint:exhaustruct
			},
			expectedResponse: delivery.NewProductResponse(&models.Product{}), //nolint:exhaustruct
		},
		{
			name:      "test full required fields of product",
			idProduct: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProduct(gomock.Any(), uint64(1), test.UserID).Return(
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

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"id": testCase.idProduct})
			req.AddCookie(&test.Cookie)
			productHandler.GetProductHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetProductList(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "offset": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsList(gomock.Any(), uint64(1), uint64(2), test.UserID).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "offset": "0"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsList(gomock.Any(), uint64(0), uint64(0), test.UserID).Return(
					[]*models.ProductInFeed{}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "10", "offset": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsList(gomock.Any(), uint64(1), uint64(10), test.UserID).Return(
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
					}, nil)
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
				}),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_list", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)
			req.AddCookie(&test.Cookie)
			productHandler.GetProductListHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetListProductOfSaler(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "last_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(2), test.UserID, true).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "last_id": "0"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(0), uint64(0), test.UserID, true).Return(
					[]*models.ProductInFeed{}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "10", "last_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(10), test.UserID, true).Return(
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
					}, nil)
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
				}),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_list_of_saler", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)
			req.AddCookie(&test.Cookie)
			productHandler.GetListProductOfSalerHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetListProductOfAnotherSaler(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "offset": "1", "saler_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(2), test.UserID, false).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "offset": "0", "saler_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(0), uint64(0), test.UserID, false).Return(
					[]*models.ProductInFeed{}, nil)
			},
			expectedResponse: delivery.NewProductListResponse(
				[]*models.ProductInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "10", "offset": "1", "saler_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetProductsOfSaler(gomock.Any(), uint64(1), uint64(10), test.UserID, false).Return(
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
					}, nil)
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

			testCase.behaviorProductService(mockProductService)

			productHandler, err := delivery.NewProductHandler("test",
				"test", "test", "test", mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_list_of_another_saler", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)

			productHandler.GetListProductOfAnotherSalerHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       any
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
  "description": "description not empty"}`)), true, uint64(1), test.UserID).Return(nil)
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
					``)), true, uint64(1), test.UserID).Return(nil)
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
  "description": "description not empty"}`)), false, uint64(1), test.UserID).Return(nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test wrong method",
			request: httptest.NewRequest(http.MethodDelete, "/api/v1/product/update", strings.NewReader(
				``)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
			},
			expectedResponse: `Method not allowed
`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&test.Cookie)
			productHandler.UpdateProductHandler(w, testCase.request)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

//nolint:dupl
func TestCloseProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryID                string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:    "test basic work",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().CloseProduct(gomock.Any(), uint64(1), test.UserID)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulCloseProduct},
			},
		},
		{
			name:    "test error in close",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().CloseProduct(gomock.Any(), uint64(1), test.UserID).Return(
					myerrors.NewErrorInternal("Test Error Internal"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name:    "test error uncorrected query param",
			queryID: "wrong type",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadFormatRequest,
				fmt.Sprintf("%s id=wrong type", utils.MessageErrWrongNumberParam)),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/product/close", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"id": testCase.queryID})
			req.AddCookie(&test.Cookie)
			productHandler.CloseProductHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

//nolint:dupl
func TestActivateProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryID                string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:    "test basic work",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().ActivateProduct(gomock.Any(), uint64(1), test.UserID)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulActivateProduct},
			},
		},
		{
			name:    "test error in internal activate",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().ActivateProduct(gomock.Any(), uint64(1), test.UserID).Return(
					myerrors.NewErrorInternal("Test Error Internal"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name:    "test error uncorrected query param",
			queryID: "wrong type",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadFormatRequest,
				fmt.Sprintf("%s id=wrong type", utils.MessageErrWrongNumberParam)),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/product/activate", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"id": testCase.queryID})
			req.AddCookie(&test.Cookie)
			productHandler.ActivateProductHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryID                string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:    "test basic work",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().DeleteProduct(gomock.Any(), uint64(1), test.UserID).Return(nil)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulDeleteProduct},
			},
		},
		{
			name:    "test error in internal",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().DeleteProduct(gomock.Any(), uint64(1), test.UserID).Return(
					myerrors.NewErrorInternal("Test Error Internal"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name:    "test error uncorrected query param",
			queryID: "wrong type",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadFormatRequest,
				fmt.Sprintf("%s id=wrong type", utils.MessageErrWrongNumberParam)),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/product/delete", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"id": testCase.queryID})
			req.AddCookie(&test.Cookie)
			productHandler.DeleteProductHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestSearchProduct(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		querySearched          string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:          "test basic work",
			querySearched: "ноутбук",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().SearchProduct(gomock.Any(), "ноутбук").Return(
					[]string{"ноутбук Acer", "ноутбук HP", "ноутбук Mac"}, nil)
			},
			expectedResponse: delivery.ProductInSearchListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []string{"ноутбук Acer", "ноутбук HP", "ноутбук Mac"},
			},
		},
		{
			name:          "test empty param",
			querySearched: "",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().SearchProduct(gomock.Any(), "").Return([]string{}, nil)
			},
			expectedResponse: delivery.ProductInSearchListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []string{},
			},
		},
		{
			name:          "test error in internal",
			querySearched: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().SearchProduct(gomock.Any(), "1").Return(nil,
					myerrors.NewErrorInternal("Test Error Internal"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
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

			productHandler, err := delivery.NewProductHandler("test",
				"test", "test", "test",
				mockProductService, mockSessionManagerClient)
			if err != nil {
				t.Fatalf("UnExpected err=%+v\n", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/search", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"searched": testCase.querySearched})
			req.AddCookie(&test.Cookie)
			productHandler.SearchProductHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetSearchProductFeed(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "offset": "0", "searched": "ноутбук"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetSearchProductFeed(gomock.Any(), "ноутбук",
					uint64(0), uint64(2), test.UserID).Return(
					[]*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}}, nil)
			},
			expectedResponse: delivery.ProductListResponse{
				Status: statuses.StatusResponseSuccessful,
				Body:   []*models.ProductInFeed{{ID: 1, Title: "Title"}, {ID: 2, Title: "Title2"}},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get_search_feed", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)
			req.AddCookie(&test.Cookie)
			productHandler.GetSearchProductFeedHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

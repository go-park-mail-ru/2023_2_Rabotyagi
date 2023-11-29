package delivery_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	mocksauth "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth/mocks"
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
		behaviorFavouriteService func(m *mocks.MockIFavouriteService)
		request                  *http.Request
		expectedResponse         responses.ResponseID
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add-to-fav", strings.NewReader(
				`{"product_id":1}`)),
			behaviorFavouriteService: func(m *mocks.MockIFavouriteService) {
				m.EXPECT().AddToFavourites(gomock.Any(), uint64(testUserID), io.NopCloser(strings.NewReader(
					`{"product_id":1}`))).Return(uint64(1), nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		//		{
		//			name: "test another product",
		//			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add", strings.NewReader(
		//				`{"saler_id":1,
		//"category_id" :1,
		//"title":"TItle",
		//"description":"  de scription",
		//"price":1232,
		//"available_count":12,
		//"city_id":1,
		//"delivery":false, "safe_deal":false}`)),
		//			behaviorProductService: func(m *mocks.MockIProductService) {
		//				m.EXPECT().AddProduct(gomock.Any(), io.NopCloser(strings.NewReader(
		//					`{"saler_id":1,
		//"category_id" :1,
		//"title":"TItle",
		//"description":"  de scription",
		//"price":1232,
		//"available_count":12,
		//"city_id":1,
		//"delivery":false, "safe_deal":false}`)), uint64(testUserID)).Return(uint64(1), nil)
		//			},
		//			expectedResponse: responses.ResponseID{
		//				Status: statuses.StatusRedirectAfterSuccessful,
		//				Body:   responses.ResponseBodyID{ID: 1},
		//			},
		//		},
		//		{
		//			name: "test product with images",
		//			request: httptest.NewRequest(http.MethodPost, "/api/v1/product/add", strings.NewReader(
		//				`{"saler_id":1,
		//"category_id" :1,
		//"title":"TItle",
		//"description":"  de scription",
		//"price":1232,
		//"available_count":12,
		//"city_id":1,
		//"delivery":false, "safe_deal":false,
		//"images": [{"url":"img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695"},
		//{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)),
		//			behaviorProductService: func(m *mocks.MockIProductService) {
		//				m.EXPECT().AddProduct(gomock.Any(), io.NopCloser(strings.NewReader(
		//					`{"saler_id":1,
		//"category_id" :1,
		//"title":"TItle",
		//"description":"  de scription",
		//"price":1232,
		//"available_count":12,
		//"city_id":1,
		//"delivery":false, "safe_deal":false,
		//"images": [{"url":"img/0b70d1440b896bf84adac5311fcd015a41590cc23fecb2750478a342918a9695"},
		//{"url":"8244c1507a772d2a9377dd95a9ce7d7eba646a62cbb865e597f58807e1"}]}`)), uint64(testUserID)).Return(uint64(1), nil)
		//			},
		//			expectedResponse: responses.ResponseID{
		//				Status: statuses.StatusRedirectAfterSuccessful,
		//				Body:   responses.ResponseBodyID{ID: 1},
		//			},
		//		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFavouriteService := mocks.NewMockIFavouriteService(ctrl)
			mockProductService := mocks.NewMockIProductService(ctrl)
			mockSessionManagerClient := mocksauth.NewMockSessionMangerClient(ctrl)

			behaviorSessionManagerClientCheck(mockSessionManagerClient)
			testCase.behaviorFavouriteService(mockFavouriteService)

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

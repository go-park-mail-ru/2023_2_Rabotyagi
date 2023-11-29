package delivery_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	mocksauth "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
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
			name:    "test basic work",
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

package delivery_test

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddPremium(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

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
				m.EXPECT().AddPremium(gomock.Any(), uint64(1), test.UserID)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulAddPremium},
			},
		},
		{
			name:    "test error in internal add",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddPremium(gomock.Any(), uint64(1), test.UserID).Return(
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
				fmt.Sprintf("%s product_id=wrong type", utils.MessageErrWrongNumberParam)),
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

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/premium/add", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"product_id": testCase.queryID})
			req.AddCookie(&test.Cookie)
			productHandler.AddPremiumHandler(w, req)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestRemovePremium(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

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
				m.EXPECT().RemovePremium(gomock.Any(), uint64(1), test.UserID)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfullyRemovePremium},
			},
		},
		{
			name:    "test error in internal remove",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().RemovePremium(gomock.Any(), uint64(1), test.UserID).Return(
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
				fmt.Sprintf("%s product_id=wrong type", utils.MessageErrWrongNumberParam)),
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

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/premium/remove", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"product_id": testCase.queryID})
			req.AddCookie(&test.Cookie)
			productHandler.RemovePremiumHandler(w, req)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

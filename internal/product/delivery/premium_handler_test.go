package delivery_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
)

func TestAddPremium(t *testing.T) {
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
			name:                   "test basic work",
			queryID:                "1",
			behaviorProductService: func(m *mocks.MockIProductService) {},
			expectedResponse:       responses.NewErrResponse(statuses.StatusInternalServer, "Ошибка на сервере"),
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

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/premium/add", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"product_id": testCase.queryID, "period": "1"})
			req.AddCookie(&test.Cookie)
			productHandler.AddPremiumHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

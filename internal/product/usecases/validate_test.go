package usecases_test

import (
	"io"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
)

//nolint:nolintlint,funlen
func TestValidatePreProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type testCase struct {
		name               string
		inputReader        io.Reader
		expectedPreProduct *models.PreProduct
		expectedError      error
	}

	testCases := [...]testCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(`{"available_count": 1,
			"category_id": 1,  "city_id": 1, "saler_id": 1,
			"title": "title", "price" : 123,
			"description": "description not empty", "delivery":false, 
			"safe_deal":false, "is_active":true, "images":[{"url":"test_url"}]}`),
			expectedPreProduct: &models.PreProduct{
				Description: "description not empty", CityID: 1,
				CategoryID: 1, Title: "title", Price: 123,
				AvailableCount: 1, SalerID: test.UserID,
				SafeDeal: false, Delivery: false, IsActive: true,
				Images: []models.Image{{URL: "test_url"}},
			},
			expectedError: nil,
		},
		{
			name:               "test error decode",
			inputReader:        strings.NewReader(`{`),
			expectedPreProduct: nil,
			expectedError:      usecases.ErrDecodePreProduct,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			preProduct, err := usecases.ValidatePartOfPreProduct(testCase.inputReader, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(preProduct, testCase.expectedPreProduct); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:nolintlint,funlen
func TestValidatePartOfPreProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type testCase struct {
		name               string
		inputReader        io.Reader
		expectedPreProduct *models.PreProduct
		expectedError      error
	}

	testCases := [...]testCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(`{
        "description": "This is a test product",
        "price": 10,
        "available_count": 5
    }`),
			expectedPreProduct: &models.PreProduct{ //nolint:exhaustruct
				Description: "This is a test product", Price: 10, AvailableCount: 5, SalerID: test.UserID,
			},
			expectedError: nil,
		},
		{
			name:               "test error decode",
			inputReader:        strings.NewReader(`{`),
			expectedPreProduct: nil,
			expectedError:      usecases.ErrDecodePreProduct,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			preProduct, err := usecases.ValidatePartOfPreProduct(testCase.inputReader, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(preProduct, testCase.expectedPreProduct); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:nolintlint,funlen
func TestValidatePreOrder(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type testCase struct {
		name             string
		inputReader      io.Reader
		expectedPreOrder *models.PreOrder
		expectedError    error
	}

	testCases := [...]testCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(`{
        "product_id":1, "count":1}`),
			expectedPreOrder: &models.PreOrder{ProductID: test.ProductID, Count: 1},
			expectedError:    nil,
		},
		{
			name:             "test error decode",
			inputReader:      strings.NewReader(`{`),
			expectedPreOrder: nil,
			expectedError:    usecases.ErrDecodePreOrder,
		},
		{
			name:             "test error validate",
			inputReader:      strings.NewReader(`{"product_id":1}`),
			expectedPreOrder: nil,
			expectedError:    usecases.ErrValidatePreOrder,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			preProduct, err := usecases.ValidatePreOrder(testCase.inputReader)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(preProduct, testCase.expectedPreOrder); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:nolintlint,funlen
func TestValidateOrderChangesCount(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type testCase struct {
		name                 string
		inputReader          io.Reader
		expectedOrderChanges *models.OrderChanges
		expectedError        error
	}

	testCases := [...]testCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(`{
        "id":1, "count":1}`),
			expectedOrderChanges: &models.OrderChanges{ID: 1, Count: 1}, //nolint:exhaustruct
			expectedError:        nil,
		},
		{
			name:                 "test error decode",
			inputReader:          strings.NewReader(`{`),
			expectedOrderChanges: nil,
			expectedError:        usecases.ErrDecodeOrderChanges,
		},
		{
			name: "test error validation required count",
			inputReader: strings.NewReader(`{
        "id":1}`),
			expectedOrderChanges: nil,
			expectedError:        usecases.ErrValidateOrderChangesCount,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			preProduct, err := usecases.ValidateOrderChangesCount(testCase.inputReader)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(preProduct, testCase.expectedOrderChanges); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

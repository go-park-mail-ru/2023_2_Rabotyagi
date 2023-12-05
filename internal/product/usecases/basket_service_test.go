package usecases_test

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"io"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"

	"go.uber.org/mock/gomock"
)

func NewBasketService(ctrl *gomock.Controller,
	behaviorBasketStorage func(m *mocks.MockIBasketStorage),
) (*usecases.BasketService, error) {
	_ = my_logger.NewNop()

	mockBasketService := mocks.NewMockIBasketStorage(ctrl)

	behaviorBasketStorage(mockBasketService)

	basketService, err := usecases.NewBasketService(mockBasketService)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return basketService, nil
}

//nolint:funlen
func TestAddOrder(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                  string
		inputReader           io.Reader
		behaviorBasketStorage func(m *mocks.MockIBasketStorage)
		expectedOrderInBasket *models.OrderInBasket
		expectedError         error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(
				`{"product_id": 1, 
					"count": 1 }`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().AddOrderInBasket(baseCtx, test.UserID, test.ProductID, uint32(1)).Return(
					&models.OrderInBasket{ //nolint:exhaustruct
						ID: 1, ProductID: test.ProductID, Count: 1, AvailableCount: 2, SalerID: 1,
					}, nil)
			},
			expectedOrderInBasket: &models.OrderInBasket{ //nolint:exhaustruct
				ID: 1, ProductID: test.ProductID, Count: 1, AvailableCount: 2, SalerID: 1,
			},
			expectedError: nil,
		},
		{
			name: "test validation error",
			inputReader: strings.NewReader(
				`{"product_id": 1}`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {},
			expectedOrderInBasket: nil,
			expectedError:         usecases.ErrValidatePreOrder,
		},
		{
			name: "test internal error",
			inputReader: strings.NewReader(
				`{"product_id": 1, 
					"count": 1 }`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().AddOrderInBasket(baseCtx, test.UserID, test.ProductID, uint32(1)).Return(
					nil, testInternalErr)
			},
			expectedOrderInBasket: nil,
			expectedError:         testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewBasketService(ctrl, testCase.behaviorBasketStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			orderInBasket, err := productService.AddOrder(baseCtx, testCase.inputReader, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(orderInBasket, testCase.expectedOrderInBasket); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:funlen
func TestGetOrderByUserID(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                  string
		behaviorBasketStorage func(m *mocks.MockIBasketStorage)
		expectedOrderInBasket []*models.OrderInBasket
		expectedError         error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().GetOrdersInBasketByUserID(baseCtx, test.UserID).Return(
					[]*models.OrderInBasket{
						{ID: 1, ProductID: test.ProductID, Count: 1, AvailableCount: 2, SalerID: 1},
					}, nil)
			},
			expectedOrderInBasket: []*models.OrderInBasket{
				{ID: 1, ProductID: test.ProductID, Count: 1, AvailableCount: 2, SalerID: 1},
			},
			expectedError: nil,
		},
		{
			name: "test internal error",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().GetOrdersInBasketByUserID(baseCtx, test.UserID).Return(
					nil, testInternalErr)
			},
			expectedOrderInBasket: nil,
			expectedError:         testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewBasketService(ctrl, testCase.behaviorBasketStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			ordersInBasket, err := productService.GetOrdersByUserID(baseCtx, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(ordersInBasket, testCase.expectedOrderInBasket); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:funlen
func TestUpdateOrderCount(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                  string
		inputReader           io.Reader
		behaviorBasketStorage func(m *mocks.MockIBasketStorage)
		expectedError         error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(
				`{"id": 1, 
					"count": 1 }`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().UpdateOrderCount(baseCtx, test.UserID, test.ProductID, uint32(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "test validation error",
			inputReader: strings.NewReader(
				`{"id": 1}`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {},
			expectedError:         usecases.ErrValidateOrderChangesCount,
		},
		{
			name: "test internal error",
			inputReader: strings.NewReader(
				`{"id": 1, 
					"count": 1 }`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().UpdateOrderCount(baseCtx, test.UserID, test.ProductID, uint32(1)).Return(testInternalErr)
			},
			expectedError: testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewBasketService(ctrl, testCase.behaviorBasketStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.UpdateOrderCount(baseCtx, testCase.inputReader, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

//nolint:funlen
func TestUpdateOrderStatus(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                  string
		inputReader           io.Reader
		behaviorBasketStorage func(m *mocks.MockIBasketStorage)
		expectedError         error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(
				`{"id": 1, 
					"status": 1 }`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().UpdateOrderStatus(baseCtx, test.UserID, test.ProductID, uint8(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "test validation error",
			inputReader: strings.NewReader(
				`{"id": 1}`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {},
			expectedError:         usecases.ErrValidateOrderChangesStatus,
		},
		{
			name: "test internal error",
			inputReader: strings.NewReader(
				`{"id": 1, 
					"status": 1 }`),
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().UpdateOrderStatus(baseCtx, test.UserID, test.ProductID, uint8(1)).Return(testInternalErr)
			},
			expectedError: testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewBasketService(ctrl, testCase.behaviorBasketStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.UpdateOrderStatus(baseCtx, testCase.inputReader, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

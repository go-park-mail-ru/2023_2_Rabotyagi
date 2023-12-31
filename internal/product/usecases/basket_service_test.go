package usecases_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
)

func NewBasketService(ctrl *gomock.Controller,
	behaviorBasketStorage func(m *mocks.MockIBasketStorage),
) (*usecases.BasketService, error) {
	_ = mylogger.NewNop()

	mockBasketService := mocks.NewMockIBasketStorage(ctrl)

	behaviorBasketStorage(mockBasketService)

	basketService, err := usecases.NewBasketService(mockBasketService)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return basketService, nil
}

func TestAddOrder(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

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

func TestGetOrderByUserID(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

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

func TestGetOrderNotInBasketByUserID(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                     string
		behaviorBasketStorage    func(m *mocks.MockIBasketStorage)
		expectedOrderNotInBasket []*models.OrderNotInBasket
		expectedError            error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().GetOrdersNotInBasketByUserID(baseCtx, test.UserID).Return(
					[]*models.OrderNotInBasket{
						{
							OrderInBasket: models.OrderInBasket{ProductID: 1, Title: "sofa"}, //nolint:exhaustruct
						},
						{
							OrderInBasket: models.OrderInBasket{ProductID: 2, Title: "laptop"}, //nolint:exhaustruct
						},
					}, nil)
			},
			expectedOrderNotInBasket: []*models.OrderNotInBasket{
				{
					OrderInBasket: models.OrderInBasket{ProductID: 1, Title: "sofa"}, //nolint:exhaustruct
				},
				{
					OrderInBasket: models.OrderInBasket{ProductID: 2, Title: "laptop"}, //nolint:exhaustruct
				},
			},
			expectedError: nil,
		},
		{
			name: "test internal error",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().GetOrdersNotInBasketByUserID(baseCtx, test.UserID).Return(
					nil, testInternalErr)
			},
			expectedOrderNotInBasket: nil,
			expectedError:            testInternalErr,
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

			ordersInBasket, err := productService.GetOrdersNotInBasketByUserID(baseCtx, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(ordersInBasket, testCase.expectedOrderNotInBasket); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

func TestGetOrderSoldByUserID(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                     string
		behaviorBasketStorage    func(m *mocks.MockIBasketStorage)
		expectedOrderNotInBasket []*models.OrderNotInBasket
		expectedError            error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().GetOrdersSoldByUserID(baseCtx, test.UserID).Return(
					[]*models.OrderNotInBasket{
						{
							OrderInBasket: models.OrderInBasket{ProductID: 1, Title: "sofa"}, //nolint:exhaustruct
						},
						{
							OrderInBasket: models.OrderInBasket{ProductID: 2, Title: "laptop"}, //nolint:exhaustruct
						},
					}, nil)
			},
			expectedOrderNotInBasket: []*models.OrderNotInBasket{
				{
					OrderInBasket: models.OrderInBasket{ProductID: 1, Title: "sofa"}, //nolint:exhaustruct
				},
				{
					OrderInBasket: models.OrderInBasket{ProductID: 2, Title: "laptop"}, //nolint:exhaustruct
				},
			},
			expectedError: nil,
		},
		{
			name: "test internal error",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().GetOrdersSoldByUserID(baseCtx, test.UserID).Return(
					nil, testInternalErr)
			},
			expectedOrderNotInBasket: nil,
			expectedError:            testInternalErr,
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

			ordersInBasket, err := productService.GetOrdersSoldByUserID(baseCtx, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(ordersInBasket, testCase.expectedOrderNotInBasket); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

func TestUpdateOrderCount(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

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

func TestUpdateOrderStatus(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

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

func TestBuyFullBasket(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type testCase struct {
		name                  string
		behaviorBasketStorage func(m *mocks.MockIBasketStorage)
		expectedError         error
	}

	testCases := [...]testCase{
		{
			name: "test basic work",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().BuyFullBasket(baseCtx, test.UserID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "test internal error",
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().BuyFullBasket(baseCtx, test.UserID).Return(testInternalErr)
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

			err = productService.BuyFullBasket(baseCtx, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

func TestDeleteOrder(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type testCase struct {
		name                  string
		inputOrderID          uint64
		behaviorBasketStorage func(m *mocks.MockIBasketStorage)
		expectedError         error
	}

	testCases := [...]testCase{
		{
			name:         "test basic work",
			inputOrderID: 1,
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().DeleteOrder(baseCtx, uint64(1), test.UserID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:         "test internal error",
			inputOrderID: 1,
			behaviorBasketStorage: func(m *mocks.MockIBasketStorage) {
				m.EXPECT().DeleteOrder(baseCtx, uint64(1), test.UserID).Return(testInternalErr)
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

			err = productService.DeleteOrder(baseCtx, testCase.inputOrderID, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

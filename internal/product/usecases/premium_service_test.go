package usecases_test

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
	"testing"
)

func NewPremiumService(ctrl *gomock.Controller,
	behaviorPremiumStorage func(m *mocks.MockIPremiumStorage),
) (*usecases.PremiumService, error) {
	_ = my_logger.NewNop()

	mockPremiumStorage := mocks.NewMockIPremiumStorage(ctrl)

	behaviorPremiumStorage(mockPremiumStorage)

	PremiumService, err := usecases.NewPremiumService(mockPremiumStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return PremiumService, nil
}

//nolint:funlen
func TestAddPremium(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		inputProductID         uint64
		behaviorPremiumStorage func(m *mocks.MockIPremiumStorage)
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name:           "test basic work",
			inputProductID: test.ProductID,
			behaviorPremiumStorage: func(m *mocks.MockIPremiumStorage) {
				m.EXPECT().AddPremium(baseCtx, test.ProductID, test.UserID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:           "test internal error",
			inputProductID: test.ProductID,
			behaviorPremiumStorage: func(m *mocks.MockIPremiumStorage) {
				m.EXPECT().AddPremium(baseCtx, test.ProductID, test.UserID).Return(testInternalErr)
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

			productService, err := NewPremiumService(ctrl, testCase.behaviorPremiumStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.AddPremium(baseCtx, testCase.inputProductID, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

//nolint:funlen
func TestRemovePremium(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		inputProductID         uint64
		behaviorPremiumStorage func(m *mocks.MockIPremiumStorage)
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name:           "test basic work",
			inputProductID: test.ProductID,
			behaviorPremiumStorage: func(m *mocks.MockIPremiumStorage) {
				m.EXPECT().RemovePremium(baseCtx, test.ProductID, test.UserID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:           "test internal error",
			inputProductID: test.ProductID,
			behaviorPremiumStorage: func(m *mocks.MockIPremiumStorage) {
				m.EXPECT().RemovePremium(baseCtx, test.ProductID, test.UserID).Return(testInternalErr)
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

			productService, err := NewPremiumService(ctrl, testCase.behaviorPremiumStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.RemovePremium(baseCtx, testCase.inputProductID, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

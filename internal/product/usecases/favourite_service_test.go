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
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"

	"go.uber.org/mock/gomock"
)

func NewFavoutiteService(ctrl *gomock.Controller,
	behaviorFavouriteStorage func(m *mocks.MockIFavouriteStorage),
) (*usecases.FavouriteService, error) {
	_ = my_logger.NewNop()

	mockFavouriteStorage := mocks.NewMockIFavouriteStorage(ctrl)

	behaviorFavouriteStorage(mockFavouriteStorage)

	favouriteService, err := usecases.NewFavouriteService(mockFavouriteStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return favouriteService, nil
}

//nolint:funlen
func TestGetUserFavourites(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                     string
		behaviorFavouriteStorage func(m *mocks.MockIFavouriteStorage)
		expectedProductInFeed    []*models.ProductInFeed
		expectedError            error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {
				m.EXPECT().GetUserFavourites(baseCtx, test.UserID).Return(
					[]*models.ProductInFeed{
						{ID: 1, Title: "test"},
					}, nil)
			},
			expectedProductInFeed: []*models.ProductInFeed{
				{ID: 1, Title: "test"},
			},
			expectedError: nil,
		},
		{
			name: "test internal error",
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {
				m.EXPECT().GetUserFavourites(baseCtx, test.UserID).Return(
					nil, testInternalErr)
			},
			expectedProductInFeed: nil,
			expectedError:         testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewFavoutiteService(ctrl, testCase.behaviorFavouriteStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			slProductInFeed, err := productService.GetUserFavourites(baseCtx, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(slProductInFeed, testCase.expectedProductInFeed); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:funlen
func TestAddToFavourites(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                     string
		inputReader              io.Reader
		behaviorFavouriteStorage func(m *mocks.MockIFavouriteStorage)
		expectedError            error
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			inputReader: strings.NewReader(`{"product_id":1}`),
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {
				m.EXPECT().AddToFavourites(baseCtx, test.UserID, test.ProductID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:                     "test error decode",
			inputReader:              strings.NewReader(`{"product_id":"wrong_type"}`),
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {},
			expectedError:            usecases.ErrDecodeProductID,
		},
		{
			name:        "test internal error",
			inputReader: strings.NewReader(`{"product_id":1}`),
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {
				m.EXPECT().AddToFavourites(baseCtx, test.UserID, test.ProductID).Return(testInternalErr)
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

			productService, err := NewFavoutiteService(ctrl, testCase.behaviorFavouriteStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.AddToFavourites(baseCtx, test.UserID, testCase.inputReader)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

//nolint:funlen
func TestDeleteFromFavourites(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                     string
		inputProductID           uint64
		behaviorFavouriteStorage func(m *mocks.MockIFavouriteStorage)
		expectedError            error
	}

	testCases := [...]TestCase{
		{
			name:           "test basic work",
			inputProductID: test.ProductID,
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {
				m.EXPECT().DeleteFromFavourites(baseCtx, test.UserID, test.ProductID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:           "test internal error",
			inputProductID: test.ProductID,
			behaviorFavouriteStorage: func(m *mocks.MockIFavouriteStorage) {
				m.EXPECT().DeleteFromFavourites(baseCtx, test.UserID, test.ProductID).Return(testInternalErr)
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

			productService, err := NewFavoutiteService(ctrl, testCase.behaviorFavouriteStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.DeleteFromFavourites(baseCtx, test.UserID, testCase.inputProductID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

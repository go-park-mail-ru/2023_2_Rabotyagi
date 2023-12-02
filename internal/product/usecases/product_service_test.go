package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	mocksfileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"io"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func NewProductService(ctrl *gomock.Controller,
	behaviorProductStorage func(m *mocks.MockIProductStorage), behaviorFileService func(m *mocksfileservice.MockFileServiceClient),
) (*usecases.ProductService, error) {
	mockProductStorage := mocks.NewMockIProductStorage(ctrl)
	mockFileService := mocksfileservice.NewMockFileServiceClient(ctrl)
	mockBasketStrorage := mocks.NewMockIBasketStorage(ctrl)
	mockFavouriteStrorage := mocks.NewMockIFavouriteStorage(ctrl)

	behaviorProductStorage(mockProductStorage)
	behaviorFileService(mockFileService)

	basketService, err := usecases.NewBasketService(mockBasketStrorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	favouriteService, err := usecases.NewFavouriteService(mockFavouriteStrorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	productService, err := usecases.NewProductService(mockProductStorage,
		basketService, favouriteService, mockFileService)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return productService, nil
}

//nolint:funlen
func TestAddProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()
	baseCtx := context.Background()

	type TestCase struct {
		name                      string
		behaviorProductStorage    func(m *mocks.MockIProductStorage)
		behaviorFileServiceClient func(m *mocksfileservice.MockFileServiceClient)
		inputReader               io.Reader
		expectedProductID         uint64
		expectedError             error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(
				`{"saler_id":1,
"category_id" :2,
"title":"adsf",
"description":"description",
"price":123,
"available_count":1,
"city_id":1,
"delivery":false, "safe_deal":false}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().AddProduct(baseCtx, &models.PreProduct{
					SalerID:        1,
					CategoryID:     2,
					Title:          "adsf",
					Description:    "description",
					Price:          123,
					AvailableCount: 1,
					CityID:         1,
					Delivery:       false,
					SafeDeal:       false,
				}).Return(test.ProductID, nil)
			},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {},
			expectedProductID:         test.ProductID,
			expectedError:             nil,
		},
		{
			name: "test bad format ",
			inputReader: strings.NewReader(
				`{`),
			behaviorProductStorage:    func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {},
			expectedProductID:         0,
			expectedError:             usecases.ErrDecodePreProduct,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewProductService(ctrl, testCase.behaviorProductStorage, testCase.behaviorFileServiceClient)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			productID, err := productService.AddProduct(baseCtx, testCase.inputReader, test.UserID)
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Failed AddProduct: err got %+v err expected: %+v", err, testCase.expectedError)
			}

			if err := utils.CompareSameType(productID, testCase.expectedProductID); err != nil {
				t.Fatalf("Failed CompareSameType %+v", err)
			}
		})
	}
}

package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	fileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service"
	mocksfileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"

	"go.uber.org/mock/gomock"
)

func NewProductService(ctrl *gomock.Controller,
	behaviorProductStorage func(m *mocks.MockIProductStorage),
	behaviorFileService func(m *mocksfileservice.MockFileServiceClient),
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
	testInternalErr := myerrors.NewErrorInternal("Test error")

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
				m.EXPECT().AddProduct(baseCtx, &models.PreProduct{ //nolint:exhaustruct
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
		{
			name: "test check urls internal error",
			inputReader: strings.NewReader(
				`{"saler_id":1,
			"category_id" :2,
			"title":"adsf",
			"description":"description",
			"price":123,
			"available_count":1,
			"city_id":1,
			"delivery":false, "safe_deal":false,
			"images": [{"url": "test_url"}]}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					nil, testInternalErr)
			},
			expectedProductID: 0,
			expectedError:     testInternalErr,
		},
		{
			name: "test addProduct internal error",
			inputReader: strings.NewReader(
				`{"saler_id":1,
			"category_id" :2,
			"title":"adsf",
			"description":"description",
			"price":123,
			"available_count":1,
			"city_id":1,
			"delivery":false, "safe_deal":false,
			"images": [{"url": "test_url"}]}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().AddProduct(baseCtx, &models.PreProduct{ //nolint:exhaustruct
					SalerID:        1,
					CategoryID:     2,
					Title:          "adsf",
					Description:    "description",
					Price:          123,
					AvailableCount: 1,
					CityID:         1,
					Delivery:       false,
					SafeDeal:       false,
					Images:         []models.Image{{URL: "test_url"}},
				}).Return(uint64(0), testInternalErr)
			},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					&fileservice.CheckedURLs{Correct: []bool{true}}, nil)
			},
			expectedProductID: 0,
			expectedError:     testInternalErr,
		},
		{
			name: "test work with images",
			inputReader: strings.NewReader(
				`{"saler_id":1,
			"category_id" :2,
			"title":"adsf",
			"description":"description",
			"price":123,
			"available_count":1,
			"city_id":1,
			"delivery":false, "safe_deal":false,
			"images": [{"url": "test_url"}]}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().AddProduct(baseCtx, &models.PreProduct{ //nolint:exhaustruct
					SalerID:        1,
					CategoryID:     2,
					Title:          "adsf",
					Description:    "description",
					Price:          123,
					AvailableCount: 1,
					CityID:         1,
					Delivery:       false,
					SafeDeal:       false,
					Images:         []models.Image{{URL: "test_url"}},
				}).Return(test.ProductID, nil)
			},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					&fileservice.CheckedURLs{Correct: []bool{true}}, nil)
			},
			expectedProductID: test.ProductID,
			expectedError:     nil,
		},
		{
			name: "test checkedURLs == nil",
			inputReader: strings.NewReader(
				`{"saler_id":1,
			"category_id" :2,
			"title":"adsf",
			"description":"description",
			"price":123,
			"available_count":1,
			"city_id":1,
			"delivery":false, "safe_deal":false,
			"images": [{"url": "test_url"}]}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					nil, nil)
			},
			expectedProductID: 0,
			expectedError:     myerrors.NewErrorInternal("checkedURLs == nil"),
		},
		{
			name: "test different len checkedURLs and requested urls",
			inputReader: strings.NewReader(
				`{"saler_id":1,
			"category_id" :2,
			"title":"adsf",
			"description":"description",
			"price":123,
			"available_count":1,
			"city_id":1,
			"delivery":false, "safe_deal":false,
			"images": [{"url": "test_url"}]}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					&fileservice.CheckedURLs{Correct: []bool{true, true}}, nil)
			},
			expectedProductID: 0,
			expectedError:     myerrors.NewErrorInternal("Different lens of checkedURLs.Correct and slImg 2 != 1"),
		},
		{
			name: "test uncorrected url",
			inputReader: strings.NewReader(
				`{"saler_id":1,
			"category_id" :2,
			"title":"adsf",
			"description":"description",
			"price":123,
			"available_count":1,
			"city_id":1,
			"delivery":false, "safe_deal":false,
			"images": [{"url": "test_url"}]}`),
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					&fileservice.CheckedURLs{Correct: []bool{false}}, nil)
			},
			expectedProductID: 0,
			expectedError:     myerrors.NewErrorBadFormatRequest("файл с урлом: test_url не найден в хранилище"),
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
				if !(err.Error() == testCase.expectedError.Error()) {
					t.Fatalf("Failed AddProduct: err got %+v err expected: %+v", err, testCase.expectedError)
				}
			}

			if err := utils.CompareSameType(productID, testCase.expectedProductID); err != nil {
				t.Fatalf("Failed CompareSameType %+v", err)
			}
		})
	}
}

//nolint:funlen
func TestGetProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *mocks.MockIProductStorage)
		inputProductID         uint64
		expectedProductID      *models.Product
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name:           "test basic work",
			inputProductID: test.ProductID,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().GetProduct(baseCtx, test.ProductID, test.UserID).Return(&models.Product{ //nolint:exhaustruct
					ID: test.ProductID, Title: "Test",
				}, nil)
			},
			expectedProductID: &models.Product{ //nolint:exhaustruct
				ID: test.ProductID, Title: "Test",
			},
			expectedError: nil,
		},
		{
			name:           "test internal error",
			inputProductID: test.ProductID,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().GetProduct(baseCtx, test.ProductID, test.UserID).Return(nil, testInternalErr)
			},
			expectedProductID: nil,
			expectedError:     testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewProductService(ctrl, testCase.behaviorProductStorage,
				func(m *mocksfileservice.MockFileServiceClient) {})
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			product, err := productService.GetProduct(baseCtx, testCase.inputProductID, test.UserID)
			if !errors.Is(err, testCase.expectedError) {
				if !(err.Error() == testCase.expectedError.Error()) {
					t.Fatalf("Failed AddProduct: err got %+v err expected: %+v", err, testCase.expectedError)
				}
			}

			if err := utils.EqualTest(product, testCase.expectedProductID); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:funlen
func TestGetProductList(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *mocks.MockIProductStorage)
		inputLastProductID     uint64
		inputCount             uint64
		expectedProductID      []*models.ProductInFeed
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name:               "test basic work",
			inputLastProductID: test.ProductID,
			inputCount:         test.CountProduct,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().GetOldProducts(baseCtx, test.ProductID, test.CountProduct, test.UserID).Return(
					[]*models.ProductInFeed{
						{ID: test.ProductID, Title: "Title"}, {ID: test.ProductID + 1, Title: "Title"},
					}, nil)
			},
			expectedProductID: []*models.ProductInFeed{
				{ID: test.ProductID, Title: "Title"}, {ID: test.ProductID + 1, Title: "Title"},
			},
			expectedError: nil,
		},
		{
			name:               "test internal error",
			inputLastProductID: test.ProductID,
			inputCount:         test.CountProduct,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().GetOldProducts(baseCtx, test.ProductID, test.CountProduct, test.UserID).Return(
					nil, testInternalErr)
			},
			expectedProductID: nil,
			expectedError:     testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewProductService(ctrl, testCase.behaviorProductStorage,
				func(m *mocksfileservice.MockFileServiceClient) {})
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			product, err := productService.GetProductsList(baseCtx,
				testCase.inputLastProductID, test.CountProduct, test.UserID)
			if !errors.Is(err, testCase.expectedError) {
				if !(err.Error() == testCase.expectedError.Error()) {
					t.Fatalf("Failed AddProduct: err got %+v err expected: %+v", err, testCase.expectedError)
				}
			}

			if err := utils.EqualTest(product, testCase.expectedProductID); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

//nolint:funlen
func TestGetProductsOfSaler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *mocks.MockIProductStorage)
		inputLastProductID     uint64
		inputCount             uint64
		expectedProductID      []*models.ProductInFeed
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name:               "test basic work",
			inputLastProductID: test.ProductID,
			inputCount:         test.CountProduct,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().GetProductsOfSaler(baseCtx, test.ProductID, test.CountProduct, test.UserID, true).Return(
					[]*models.ProductInFeed{
						{ID: test.ProductID, Title: "Title"}, {ID: test.ProductID + 1, Title: "Title"},
					}, nil)
			},
			expectedProductID: []*models.ProductInFeed{
				{ID: test.ProductID, Title: "Title"}, {ID: test.ProductID + 1, Title: "Title"},
			},
			expectedError: nil,
		},
		{
			name:               "test internal error",
			inputLastProductID: test.ProductID,
			inputCount:         test.CountProduct,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().GetProductsOfSaler(baseCtx, test.ProductID, test.CountProduct, test.UserID, true).Return(
					nil, testInternalErr)
			},
			expectedProductID: nil,
			expectedError:     testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewProductService(ctrl, testCase.behaviorProductStorage,
				func(m *mocksfileservice.MockFileServiceClient) {})
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			product, err := productService.GetProductsOfSaler(baseCtx,
				testCase.inputLastProductID, test.CountProduct, test.UserID, true)
			if !errors.Is(err, testCase.expectedError) {
				if !(err.Error() == testCase.expectedError.Error()) {
					t.Fatalf("Failed AddProduct: err got %+v err expected: %+v", err, testCase.expectedError)
				}
			}

			if err := utils.EqualTest(product, testCase.expectedProductID); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

func generateString(lenStr int) string {
	templateStr := "aaaaaaaaaa"
	result := ""

	for i := 0; i < lenStr/10; i++ {
		result += templateStr
	}

	for len(result) < lenStr {
		result += "a"
	}

	return result
}

//nolint:funlen
func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                      string
		inputReader               io.Reader
		inputPartialUpdate        bool
		inputProductID            uint64
		behaviorProductStorage    func(m *mocks.MockIProductStorage)
		behaviorFileServiceClient func(m *mocksfileservice.MockFileServiceClient)
		expectedError             error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: io.NopCloser(strings.NewReader(`{"available_count": 1,
  "description": "description empty",
  "title": "Product",
  "images": [
    {
      "url": "test_url"
    }]
  }`)),
			inputPartialUpdate: true,
			inputProductID:     test.ProductID,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().UpdateProduct(baseCtx, test.ProductID,
					map[string]any{
						"available_count": uint32(1), "delivery": false,
						"description": "description empty", "images": []models.Image{{URL: "test_url"}},
						"is_active": false, "safe_deal": false, "saler_id": uint64(1), "title": "Product",
					}).Return(nil)
			},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					&fileservice.CheckedURLs{Correct: []bool{true}}, nil)
			},
			expectedError: nil,
		},
		{
			name: "test internal error update",
			inputReader: io.NopCloser(strings.NewReader(`{"available_count": 1,
  "description": "description empty",
  "title": "Product",
  "images": [
    {
      "url": "test_url"
    }]
  }`)),
			inputPartialUpdate: true,
			inputProductID:     test.ProductID,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {
				m.EXPECT().UpdateProduct(baseCtx, test.ProductID,
					map[string]any{
						"available_count": uint32(1), "delivery": false,
						"description": "description empty", "images": []models.Image{{URL: "test_url"}},
						"is_active": false, "safe_deal": false, "saler_id": uint64(1), "title": "Product",
					}).Return(testInternalErr)
			},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					&fileservice.CheckedURLs{Correct: []bool{true}}, nil)
			},
			expectedError: testInternalErr,
		},
		{
			name: "test internal error check ",
			inputReader: io.NopCloser(strings.NewReader(`{"available_count": 1,
  "description": "description empty",
  "title": "Product",
  "images": [
    {
      "url": "test_url"
    }]
  }`)),
			inputPartialUpdate:     true,
			inputProductID:         test.ProductID,
			behaviorProductStorage: func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {
				m.EXPECT().Check(baseCtx, &fileservice.ImgURLs{Url: []string{"test_url"}}).Return(
					nil, testInternalErr)
			},
			expectedError: testInternalErr,
		},
		{
			name: "test validation error long title",
			inputReader: io.NopCloser(strings.NewReader(fmt.Sprintf(`{"available_count": 1,
  "description": "1",
  "title": "%s",
  "images": [
    {
      "url": "test_url"
    }]
  }`, generateString(257)))),
			inputPartialUpdate:        true,
			inputProductID:            test.ProductID,
			behaviorProductStorage:    func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {},
			expectedError:             fmt.Errorf("в поле title ошибка: Заголовок должен быть длинной от 1 до 256 символов"),
		},
		{
			name:                      "test validation error on PUT update",
			inputReader:               io.NopCloser(strings.NewReader(`{"available_count": 1}`)),
			inputPartialUpdate:        false,
			inputProductID:            test.ProductID,
			behaviorProductStorage:    func(m *mocks.MockIProductStorage) {},
			behaviorFileServiceClient: func(m *mocksfileservice.MockFileServiceClient) {},
			expectedError: fmt.Errorf("category_id: non zero value required;city_id:" +
				" non zero value required;description: non zero value required;price:" +
				" non zero value required;title: non zero value required"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewProductService(ctrl, testCase.behaviorProductStorage,
				testCase.behaviorFileServiceClient)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.UpdateProduct(baseCtx, testCase.inputReader,
				testCase.inputPartialUpdate, testCase.inputProductID, test.UserID)
			if !errors.Is(err, testCase.expectedError) {
				if !(err.Error() == testCase.expectedError.Error()) {
					t.Fatalf("Failed AddProduct: err got %+v err expected: %+v", err, testCase.expectedError)
				}
			}
		})
	}
}

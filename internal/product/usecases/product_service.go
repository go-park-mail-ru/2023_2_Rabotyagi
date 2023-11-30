package usecases

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	fileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/microcosm-cc/bluemonday"
)

var _ IProductStorage = (*productrepo.ProductStorage)(nil)

type IProductStorage interface {
	AddProduct(ctx context.Context, preProduct *models.PreProduct) (uint64, error)
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error)
	GetOldProducts(ctx context.Context, lastProductID uint64, count uint64, userID uint64) ([]*models.ProductInFeed, error)
	GetProductsOfSaler(ctx context.Context, lastProductID uint64,
		count uint64, userID uint64, isMy bool) ([]*models.ProductInFeed, error)
	UpdateProduct(ctx context.Context, productID uint64, updateFields map[string]interface{}) error
	CloseProduct(ctx context.Context, productID uint64, userID uint64) error
	ActivateProduct(ctx context.Context, productID uint64, userID uint64) error
	DeleteProduct(ctx context.Context, productID uint64, userID uint64) error
	SearchProduct(ctx context.Context, searchInput string) ([]string, error)
	GetSearchProductFeed(ctx context.Context,
		searchInput string, lastNumber uint64, limit uint64, userID uint64,
	) ([]*models.ProductInFeed, error)
	IBasketStorage
	IFavouriteStorage
}

type ProductService struct {
	FavouriteService
	BasketService
	fileServiceClient fileservice.FileServiceClient
	storage           IProductStorage
	logger            *my_logger.MyLogger
}

func NewProductService(productStorage IProductStorage, basketService BasketService,
	favouriteService FavouriteService, fileServiceClient fileservice.FileServiceClient,
) (*ProductService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductService{
		FavouriteService:  favouriteService,
		BasketService:     basketService,
		fileServiceClient: fileServiceClient,
		storage:           productStorage,
		logger:            logger,
	}, nil
}

func (p *ProductService) checkCorrectnessUrlsImg(ctx context.Context, slImg []models.Image) error {
	logger := p.logger.LogReqID(ctx)

	checkedURLs, err := p.fileServiceClient.Check(
		ctx, &fileservice.ImgURLs{Url: convertImagesToSl(slImg)})
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if checkedURLs == nil {
		err := myerrors.NewErrorInternal("checkedURLs == nil")
		logger.Errorln(err)

		return err
	}

	if len(checkedURLs.Correct) != len(slImg) {
		err := myerrors.NewErrorInternal(
			"Different lens of checkedURLs.Correct and slImg %d != %d",
			len(checkedURLs.Correct), len(slImg))
		logger.Errorln(err)

		return err
	}

	messageUnCorrect := ""

	for i, urlCorrect := range checkedURLs.Correct {
		if !urlCorrect {
			messageUnCorrect += fmt.Sprintf("файл с урлом: %s не найден в хранилище\n", slImg[i].URL)
		}
	}

	if messageUnCorrect != "" {
		return myerrors.NewErrorBadContentRequest(messageUnCorrect)
	}

	return nil
}

func (p *ProductService) AddProduct(ctx context.Context, r io.Reader, userID uint64) (uint64, error) {
	preProduct, err := ValidatePreProduct(r, userID)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = p.checkCorrectnessUrlsImg(ctx, preProduct.Images)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	productID, err := p.storage.AddProduct(ctx, preProduct)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return productID, nil
}

func (p *ProductService) GetProduct(ctx context.Context,
	productID uint64, userID uint64,
) (*models.Product, error) {
	product, err := p.storage.GetProduct(ctx, productID, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	product.Sanitize()

	return product, nil
}

func (p *ProductService) GetProductsList(ctx context.Context,
	lastProductID uint64, count uint64, userID uint64,
) ([]*models.ProductInFeed, error) {
	products, err := p.storage.GetOldProducts(ctx, lastProductID, count, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range products {
		product.Sanitize()
	}

	return products, nil
}

func (p *ProductService) GetProductsOfSaler(ctx context.Context,
	lastProductID uint64, count uint64, userID uint64, isMy bool,
) ([]*models.ProductInFeed, error) {
	products, err := p.storage.GetProductsOfSaler(ctx, lastProductID, count, userID, isMy)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range products {
		product.Sanitize()
	}

	return products, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context,
	r io.Reader, isPartialUpdate bool, productID uint64, userAuthID uint64,
) error {
	var preProduct *models.PreProduct

	var err error

	if isPartialUpdate {
		preProduct, err = ValidatePartOfPreProduct(r, userAuthID)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}
	} else {
		preProduct, err = ValidatePreProduct(r, userAuthID)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	err = p.checkCorrectnessUrlsImg(ctx, preProduct.Images)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	updateFieldsMap := utils.StructToMap(preProduct)

	err = p.storage.UpdateProduct(ctx, productID, updateFieldsMap)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductService) CloseProduct(ctx context.Context, productID uint64, userID uint64) error {
	err := p.storage.CloseProduct(ctx, productID, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductService) ActivateProduct(ctx context.Context, productID uint64, userID uint64) error {
	err := p.storage.ActivateProduct(ctx, productID, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, productID uint64, userID uint64) error {
	err := p.storage.DeleteProduct(ctx, productID, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (p *ProductService) SearchProduct(ctx context.Context, searchInput string) ([]string, error) {
	regex := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9\s]+`)
	searchInput = regex.ReplaceAllString(searchInput, "")
	regex = regexp.MustCompile(`\s+`)
	searchInput = regex.ReplaceAllString(searchInput, " ")

	searchInput = strings.TrimSpace(searchInput)

	products, err := p.storage.SearchProduct(ctx, searchInput)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	sanitizer := bluemonday.UGCPolicy()

	for _, product := range products {
		product = sanitizer.Sanitize(product)
	}

	return products, nil
}

func (p *ProductService) GetSearchProductFeed(ctx context.Context,
	searchInput string, lastNumber uint64, limit uint64, userID uint64,
) ([]*models.ProductInFeed, error) {
	regex := regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9\s]+`)
	searchInput = regex.ReplaceAllString(searchInput, "")
	regex = regexp.MustCompile(`\s+`)
	searchInput = regex.ReplaceAllString(searchInput, " ")

	searchInput = strings.TrimSpace(searchInput)

	products, err := p.storage.GetSearchProductFeed(ctx, searchInput, lastNumber, limit, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, product := range products {
		product.Sanitize()
	}

	return products, nil
}

func convertImagesToSl(images []models.Image) []string {
	result := make([]string, len(images))

	for i, image := range images {
		result[i] = image.URL
	}

	return result
}
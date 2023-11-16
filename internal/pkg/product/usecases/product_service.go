package usecases

import (
	"context"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"

	"go.uber.org/zap"
)

var ErrUserPermissionsChange = myerrors.NewError("Вы не можете изменить чужое объявление")

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
	IBasketStorage
}

type ProductService struct {
	BasketService
	storage IProductStorage
	logger  *zap.SugaredLogger
}

func NewProductService(productStorage IProductStorage, basketService BasketService) (*ProductService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &ProductService{BasketService: basketService, storage: productStorage, logger: logger}, nil
}

func (p *ProductService) AddProduct(ctx context.Context, r io.Reader) (uint64, error) {
	preProduct, err := ValidatePreProduct(r)
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
		preProduct, err = ValidatePartOfPreProduct(r)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}
	} else {
		preProduct, err = ValidatePreProduct(r)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	if preProduct.SalerID != userAuthID {
		p.logger.Errorln(ErrUserPermissionsChange)

		return fmt.Errorf(myerrors.ErrTemplate, ErrUserPermissionsChange)
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

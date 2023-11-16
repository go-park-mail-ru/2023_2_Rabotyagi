package usecases

import (
	"context"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"

	"go.uber.org/zap"
)

var _ IProductStorage = (*productrepo.ProductStorage)(nil)

type IBasketStorage interface {
	AddOrderInBasket(ctx context.Context, userID uint64, productID uint64, count uint32) (*models.OrderInBasket, error)
	GetOrdersInBasketByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error)
	UpdateOrderCount(ctx context.Context, userID uint64, orderID uint64, newCount uint32) error
	UpdateOrderStatus(ctx context.Context, userID uint64, orderID uint64, newStatus uint8) error
	BuyFullBasket(ctx context.Context, userID uint64) error
	DeleteOrder(ctx context.Context, orderID uint64, ownerID uint64) error
}

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
	storage IProductStorage
	logger  *zap.SugaredLogger
}

func NewProductService(productStorage IProductStorage, logger *zap.SugaredLogger) *ProductService {
	return &ProductService{storage: productStorage, logger: logger}
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

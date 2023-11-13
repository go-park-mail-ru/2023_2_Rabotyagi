package usecases

import (
	"context"
	productrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
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
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error)
	GetNewProducts(ctx context.Context, lastProductID uint64, count uint64, userID uint64) ([]*models.ProductInFeed, error)
	GetProductsOfSaler(ctx context.Context, lastProductID uint64,
		count uint64, userID uint64, isMy bool) ([]*models.ProductInFeed, error)
	AddProduct(ctx context.Context, preProduct *models.PreProduct) (uint64, error)
	UpdateProduct(ctx context.Context, productID uint64, updateFields map[string]interface{}) error
	CloseProduct(ctx context.Context, productID uint64, userID uint64) error
	ActivateProduct(ctx context.Context, productID uint64, userID uint64) error
	DeleteProduct(ctx context.Context, productID uint64, userID uint64) error
	IBasketStorage
}

package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

type IProductStorage interface {
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error)
	GetNewProducts(ctx context.Context, lastProductID uint64, count uint64, userID uint64) ([]*models.ProductInFeed, error)
	AddProduct(ctx context.Context, preProduct *models.PreProduct) (uint64, error)

	GetOrdersInBasketByUserID(ctx context.Context, userID uint64) ([]*models.OrderInBasket, error)
	UpdateOrderCount(ctx context.Context, orderID uint64, newCount uint32) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint64, newStatus uint8) (*models.Order, error)
}

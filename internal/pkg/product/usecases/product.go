package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

type IProductStorage interface {
	GetProduct(ctx context.Context, productID uint64, userID uint64) (*models.Product, error)
	GetNProducts(ctx context.Context) ([]*models.Product, error)
	AddProduct(ctx context.Context, preProduct *models.PreProduct) (uint64, error)
}

package category

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
)

type ICategoryStorage interface {
	GetFullCategories(ctx context.Context) ([]*models.Category, error)
	SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error)
}

type ICategoryService interface {
	GetFullCategories(ctx context.Context) ([]*models.Category, error)
	SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error)
}

type Tables interface {
	Category() string
}

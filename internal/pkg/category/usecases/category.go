package usecases

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

type ICategoryStorage interface {
	GetFullCategories(ctx context.Context) ([]*models.Category, error)
}

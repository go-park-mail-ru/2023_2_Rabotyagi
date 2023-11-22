package usecases

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/myerrors"
	"go.uber.org/zap"
)

type ICategoryStorage interface {
	GetFullCategories(ctx context.Context) ([]*models.Category, error)
}

type CategoryService struct {
	storage ICategoryStorage
	logger  *zap.SugaredLogger
}

func NewCategoryService(categoryStorage ICategoryStorage) (*CategoryService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &CategoryService{storage: categoryStorage, logger: logger}, nil
}

func (c *CategoryService) GetFullCategories(ctx context.Context) ([]*models.Category, error) {
	categories, err := c.storage.GetFullCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, category := range categories {
		category.Sanitize()
	}

	return categories, nil
}

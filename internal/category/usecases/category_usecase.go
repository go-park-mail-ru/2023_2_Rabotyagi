package usecases

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

type ICategoryStorage interface {
	GetFullCategories(ctx context.Context) ([]*models.Category, error)
	SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error)
}

type CategoryService struct {
	storage ICategoryStorage
	logger  *mylogger.MyLogger
}

func NewCategoryService(categoryStorage ICategoryStorage) (*CategoryService, error) {
	logger, err := mylogger.Get()
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

func (c *CategoryService) SearchCategory(ctx context.Context, searchInput string) ([]*models.Category, error) {
	categories, err := c.storage.SearchCategory(ctx, searchInput)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, category := range categories {
		category.Sanitize()
	}

	return categories, nil
}

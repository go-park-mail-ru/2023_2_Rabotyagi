package usecases

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/myerrors"
	"go.uber.org/zap"
)

type ICityStorage interface {
	GetFullCities(ctx context.Context) ([]*models.City, error)
	SearchCity(ctx context.Context, searchInput string) ([]*models.City, error)
}

type CityService struct {
	storage ICityStorage
	logger  *zap.SugaredLogger
}

func NewCityService(cityStorage ICityStorage) (*CityService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &CityService{storage: cityStorage, logger: logger}, nil
}

func (c *CityService) GetFullCities(ctx context.Context) ([]*models.City, error) {
	cities, err := c.storage.GetFullCities(ctx)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, city := range cities {
		city.Sanitize()
	}

	return cities, nil
}

func (c *CityService) SearchCity(ctx context.Context, searchInput string) ([]*models.City, error) {
	cities, err := c.storage.SearchCity(ctx, searchInput)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, city := range cities {
		city.Sanitize()
	}

	return cities, nil
}

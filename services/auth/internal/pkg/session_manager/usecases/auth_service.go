package usecases

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/repository"
	"go.uber.org/zap"
)

var _ IAuthStorage = (*repository.AuthStorage)(nil)

type IAuthStorage interface {
	AddUser(ctx context.Context, preUser *models.User) (*models.User, error)
	GetUser(ctx context.Context, email string, password string) (*models.User, error)
}

type AuthService struct {
	storage IAuthStorage
	logger  *zap.SugaredLogger
}

func NewAuthService(authStorage IAuthStorage) (*AuthService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &AuthService{storage: authStorage, logger: logger}, nil
}

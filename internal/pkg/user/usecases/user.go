package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

type IUserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	//CreateUser(ctx context.Context, preUser *models.UserWithoutID) error
	IsUserExist(ctx context.Context, email string, phone string) (bool, error)
	AddUser(ctx context.Context, preUser *models.UserWithoutID) (*models.User, error)
}

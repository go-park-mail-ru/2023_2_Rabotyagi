package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

type IUserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	AddUser(ctx context.Context, preUser *models.UserWithoutID) (*models.User, error)
	UpdateUser(ctx context.Context, userID uint64, updateData map[string]interface{}) error
	IsEmailBusy(ctx context.Context, email string) (bool, error) // TODO maybe unuseful in outside
	IsPhoneBusy(ctx context.Context, phone string) (bool, error) // TODO maybe unuseful in outside
}

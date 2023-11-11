package usecases

import (
	"context"

	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

var _ IUserStorage = (*userrepo.UserStorage)(nil)

type IUserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error) // TODO maybe unuseful
	GetUserWithoutPasswordByID(ctx context.Context, id uint64) (*models.UserWithoutPassword, error)
	AddUser(ctx context.Context, preUser *models.UserWithoutID) (*models.User, error)
	GetUser(ctx context.Context, email string, password string) (*models.UserWithoutPassword, error)
	UpdateUser(ctx context.Context, userID uint64, updateData map[string]interface{}) (*models.UserWithoutPassword, error)
	IsEmailBusy(ctx context.Context, email string) (bool, error) // TODO maybe unuseful in outside
	IsPhoneBusy(ctx context.Context, phone string) (bool, error) // TODO maybe unuseful in outside
}

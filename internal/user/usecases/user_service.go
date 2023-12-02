package usecases

import (
	"context"
	"fmt"
	"io"

	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

var ErrWrongUserID = myerrors.NewErrorBadFormatRequest("Попытка изменить данные другого пользователя")

var _ IUserStorage = (*userrepo.UserStorage)(nil)

type IUserStorage interface {
	GetUserWithoutPasswordByID(ctx context.Context, id uint64) (*models.UserWithoutPassword, error)
	UpdateUser(ctx context.Context, userID uint64, updateData map[string]interface{}) (*models.UserWithoutPassword, error)
}

type UserService struct {
	storage IUserStorage
	logger  *my_logger.MyLogger
}

func NewUserService(userStorage IUserStorage) (*UserService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &UserService{storage: userStorage, logger: logger}, nil
}

func (u *UserService) GetUserWithoutPasswordByID(ctx context.Context,
	userID uint64,
) (*models.UserWithoutPassword, error) {
	user, err := u.storage.GetUserWithoutPasswordByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Sanitize()

	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, r io.Reader,
	isPartialUpdate bool, userID uint64,
) (*models.UserWithoutPassword, error) {
	logger := u.logger.LogReqID(ctx)

	var userWithoutPassword *models.UserWithoutPassword

	var err error

	if isPartialUpdate {
		userWithoutPassword, err = ValidatePartOfUserWithoutPassword(r)
		if err != nil {
			return nil, fmt.Errorf(myerrors.ErrTemplate, err)
		}
	} else {
		userWithoutPassword, err = ValidateUserWithoutPassword(r)
		if err != nil {
			return nil, fmt.Errorf(myerrors.ErrTemplate, err)
		}
	}

	if userWithoutPassword.ID != userID && userWithoutPassword.ID != 0 {
		logger.Errorln(ErrWrongUserID)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrWrongUserID)
	}

	updateDataMap := utils.StructToMap(userWithoutPassword)

	delete(updateDataMap, "ID")

	updatedUser, err := u.storage.UpdateUser(ctx, userID, updateDataMap)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	updatedUser.Sanitize()

	return updatedUser, nil
}

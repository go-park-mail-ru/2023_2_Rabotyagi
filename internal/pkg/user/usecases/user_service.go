package usecases

import (
	"context"
	"fmt"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
	"go.uber.org/zap"
	"io"
	"strconv"

	userrepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

var ErrWrongUserID = myerrors.NewError("Попытка изменить данные другого пользователя")

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

type UserService struct {
	storage IUserStorage
	logger  *zap.SugaredLogger
}

func NewUserService(userStorage IUserStorage) (*UserService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &UserService{storage: userStorage, logger: logger}, nil
}

func (u *UserService) AddUser(ctx context.Context, r io.Reader) (*models.User, error) {
	userWithoutID, err := ValidateUserWithoutID(r)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutID.Password, err = utils.HashPass(userWithoutID.Password)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user, err := u.storage.AddUser(ctx, userWithoutID)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return user, nil
}

func (u *UserService) GetUser(ctx context.Context, email string, password string) (*models.UserWithoutPassword, error) {
	userWithoutID, err := ValidateUserCredentials(email, password)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user, err := u.storage.GetUser(ctx, userWithoutID.Email, userWithoutID.Password)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Sanitize()

	return user, nil
}

func (u *UserService) GetUserWithoutPasswordByID(ctx context.Context,
	userIDStr string,
) (*models.UserWithoutPassword, error) {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

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
		u.logger.Errorln(ErrWrongUserID)

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

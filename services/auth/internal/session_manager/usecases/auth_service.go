package usecases

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/session_manager/repository"
)

var ErrWrongCredentials = myerrors.NewErrorBadContentRequest("Некорректный логин или пароль")

var _ IAuthStorage = (*repository.AuthStorage)(nil)

type IAuthStorage interface {
	AddUser(ctx context.Context, email string, password string) (*models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	storage IAuthStorage
	logger  *my_logger.MyLogger
}

func NewAuthService(authStorage IAuthStorage) (*AuthService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &AuthService{storage: authStorage, logger: logger}, nil
}

func (a *AuthService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	logger := a.logger.LogReqID(ctx)

	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	hashPass, err := hex.DecodeString(user.Password)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if !utils.ComparePassAndHash(hashPass, password) {
		return "", ErrWrongCredentials
	}

	jwtPayload := jwt.UserJwtPayload{} //nolint:exhaustruct

	jwtPayload.UserID = user.ID
	jwtPayload.Email = user.Email
	jwtPayload.Expire = time.Now().Add(jwt.TimeTokenLife).Unix()

	rawJwt, err := jwt.GenerateJwtToken(&jwtPayload, jwt.GetSecret())
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return rawJwt, nil
}

func (a *AuthService) AddUser(ctx context.Context, email string, password string) (string, error) {
	logger := a.logger.LogReqID(ctx)

	password, err := utils.HashPass(password)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user, err := a.storage.AddUser(ctx, email, password)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	jwtPayload := jwt.UserJwtPayload{}

	jwtPayload.UserID = user.ID
	jwtPayload.Email = user.Email
	jwtPayload.Expire = time.Now().Add(jwt.TimeTokenLife).Unix()

	rawJwt, err := jwt.GenerateJwtToken(&jwtPayload, jwt.GetSecret())
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return rawJwt, nil
}

func (a *AuthService) Delete(ctx context.Context, rawJwt string) (string, error) {
	logger := a.logger.LogReqID(ctx)

	jwtPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.GetSecret())
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	jwtPayload.Expire = time.Now().Unix()

	newRawJwt, err := jwt.GenerateJwtToken(jwtPayload, jwt.GetSecret())
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return newRawJwt, nil
}

func (a *AuthService) Check(_ context.Context, rawJwt string) (uint64, error) {
	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.GetSecret())
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return userPayload.UserID, nil
}

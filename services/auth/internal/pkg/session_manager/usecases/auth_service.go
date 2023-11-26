package usecases

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/pkg/session_manager/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var (
	ErrWrongCredentials = myerrors.NewErrorBadContentRequest("Некорректный логин или пароль")
)

var _ IAuthStorage = (*repository.AuthStorage)(nil)

type IAuthStorage interface {
	AddUser(ctx context.Context, email string, password string) (*models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
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

func (a *AuthService) GetUserRawJWT(ctx context.Context, email string, password string) (string, error) {
	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		a.logger.Errorln(err)

		return "", status.Errorf(codes.Internal, "can`t get user")
	}

	hashPass, err := hex.DecodeString(user.Password)
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if !utils.ComparePassAndHash(hashPass, password) {
		return "", ErrWrongCredentials
	}

	jwtPayload := jwt.UserJwtPayload{}

	jwtPayload.UserID = user.ID
	jwtPayload.Email = user.Email
	jwtPayload.Expire = time.Now().Add(jwt.TimeTokenLife).Unix()

	rawJwt, err := jwt.GenerateJwtToken(&jwtPayload, jwt.GetSecret())
	if err != nil {
		a.logger.Errorln(err)

		return "", status.Errorf(codes.Internal, "can`t get user")
	}

	return rawJwt, nil
}

func (a *AuthService) AddUser(ctx context.Context, email string, password string) (string, error) {
	password, err := utils.HashPass(password)
	if err != nil {
		a.logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user, err := a.storage.AddUser(ctx, email, password)
	if err != nil {
		a.logger.Errorln(err)

		return "", status.Errorf(codes.Internal, "can`t add user")
	}

	jwtPayload := jwt.UserJwtPayload{}

	jwtPayload.UserID = user.ID
	jwtPayload.Email = user.Email
	jwtPayload.Expire = time.Now().Add(jwt.TimeTokenLife).Unix()

	rawJwt, err := jwt.GenerateJwtToken(&jwtPayload, jwt.GetSecret())
	if err != nil {
		a.logger.Errorln(err)

		return "", status.Errorf(codes.Internal, "can`t add user")
	}

	return rawJwt, nil
}

func (a *AuthService) Delete(ctx context.Context, rawJwt string) (string, error) {
	jwtPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.GetSecret())
	if err != nil {
		a.logger.Errorln(err)

		return "", status.Errorf(codes.InvalidArgument, "get incorrect access_token in sessionUser")
	}

	jwtPayload.Expire = time.Now().Unix()

	newRawJwt, err := jwt.GenerateJwtToken(jwtPayload, jwt.GetSecret())
	if err != nil {
		a.logger.Errorln(err)

		return "", status.Errorf(codes.Internal, "can`t delete user")
	}

	return newRawJwt, nil
}

func (a *AuthService) Check(ctx context.Context, rawJwt string) bool {
	_, err := jwt.NewUserJwtPayload(rawJwt, jwt.GetSecret())
	if err != nil {
		return false
	}

	return true
}

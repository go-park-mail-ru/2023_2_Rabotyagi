package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	GetUserByEmail = "SELECT id, email, phone, name, pass, birthday FROM public.\"user\" WHERE email=$1"
	CreateUser     = "INSERT INTO public.\"user\" (email, phone, name, pass, birthday) VALUES ($1, $2, $3, $4, $5)"
	IsUserExist    = "SELECT id FORM public.\"user\" WHERE email=$1 AND phone=$2"
	// UserNotExist = "SELECT id FORM public.\"user\" WHERE email=$1 AND phone=$2"
)

var (
	ErrExecuting = errors.NewError("error while executing")
	ErrParceRow  = errors.NewError("parcing row error")
)

type IUserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, preUser *models.UserWithoutID) error
	IsUserExist(ctx context.Context, email string) bool
	// UserNotExist(ctx context.Context, email string, phone string) bool
}

type UserStorage struct {
	Pool *pgxpool.Pool
}

func NewUserStorage(Pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		Pool: Pool,
	}
}

func (u *UserStorage) GetUserById(ctx context.Context, email string) (*models.User, error) {
	userLine := u.Pool.QueryRow(ctx, GetUserByEmail, email)

	user := models.User{
		Email: email,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Pass, &user.Birthday); err != nil {
		return nil, fmt.Errorf("w", ErrParceRow)
	}

	return &user, nil
}

func (u *UserStorage) IsUserExist(ctx context.Context, email string, phone string) (bool, error) {
	userIdRow := u.Pool.QueryRow(ctx, IsUserExist, email, phone)
	var userId string

	if err := userIdRow.Scan(userId); err != nil {
		return false, fmt.Errorf("w", ErrParceRow)
	}

	if userId != "" {
		return true, nil
	}
	return false, nil
}

func (u *UserStorage) CreateUser(ctx context.Context, preUser *models.UserWithoutID) error {
	_, err := u.Pool.Exec(ctx, CreateUser, preUser.Email, preUser.Name, preUser.Name, preUser.Pass, preUser.Phone)
	if err != nil {
		return fmt.Errorf("w", ErrExecuting)
	}

	return nil
}

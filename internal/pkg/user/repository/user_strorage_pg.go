package repository

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/jackc/pgx/v5"
)

const (
	GetUserById = "SELECT id, email, phone, name, pass, birthday FROM public.\"user\" WHERE id=$1"
	CreateUser  = "INSERT INTO public.\"user\" (email, phone, name, pass, birthday) VALUES ($1, $2, $3, $4, $5)"
	IsUserExist = "\"SELECT id FROM public.\"user\" WHERE email=$1\""
)

//var (
//	ErrUserAlreadyExist = errors.NewError("user already exist")
//	ErrUserNotExist     = errors.NewError("user not exist")
//)

type IUserStorage interface {
	GetUserById(userID uint64) (*models.User, error)
	CreateUser(user *models.PreUser) error
	IsUserExist(email string) bool
}

type UserStorage struct {
	db *pgx.Conn
}

func NewUserStorage(db *pgx.Conn) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) GetUserById(ctx context.Context, id uint64) (*models.User, error) {
	userLine := u.db.QueryRow(ctx, GetUserById, id)

	user := models.User{
		ID: id,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Pass, &user.Birthday); err != nil {
		err = fmt.Errorf("error happened in row.Scan: %w", err)

		return nil, err
	}

	return &user, nil
}

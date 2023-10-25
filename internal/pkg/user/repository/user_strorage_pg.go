package repository

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/jackc/pgx/v5"

	"log"
)

const (
	GetUserById = "SELECT id, email, phone, name, pass, birthday FROM public.\"user\" WHERE id=$1"
	CreateUser  = "INSERT INTO public.\"user\" (email, phone, name, pass, birthday) VALUES ($1, $2, $3, $4, $5)"
	IsUserExist = "SELECT id FORM public.\"user\" WHERE email=$1 AND phone=$2"
	// UserNotExist = "SELECT id FORM public.\"user\" WHERE email=$1 AND phone=$2"
)

//var (
//	ErrUserAlreadyExist = errors.NewError("user already exist")
//	ErrUserNotExist     = errors.NewError("user not exist")
//)

type IUserStorage interface {
	GetUserById(ctx context.Context, userID uint64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.PreUser) error
	IsUserExist(ctx context.Context, email string) bool
	// UserNotExist(ctx context.Context, email string, phone string) bool
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
		log.Printf("%v\n", err)

		return nil, err
	}

	return &user, nil
}

func (u *UserStorage) IsUserExist(ctx context.Context, email string, phone string) (bool, error) {
	userIdRow := u.db.QueryRow(ctx, IsUserExist, email, phone)
	var userId string

	if err := userIdRow.Scan(userId); err != nil {
		log.Printf("%v\n", err)

		return false, err
	}

	if userId != "" {
		return true, nil
	}
	return false, nil
}

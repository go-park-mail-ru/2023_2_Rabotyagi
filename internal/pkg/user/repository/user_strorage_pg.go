package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	GetUserByEmail = `SELECT id, email, phone, name, pass, birthday FROM public."user" WHERE email=$1`
	GetUserByID    = `SELECT id, email, phone, name, pass, birthday FROM public."user" WHERE id=$1`
	CreateUser     = `INSERT INTO public."user" (email, phone, name, pass, birthday) VALUES ($1, $2, $3, $4, $5)`
	IsUserExist    = `SELECT id FROM public."user" WHERE email=$1 OR phone=$2`
)

type IUserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	CreateUser(ctx context.Context, preUser *models.UserWithoutID) error
	IsEmailBusy(ctx context.Context, email string) bool
	IsPhoneBusy(ctx context.Context, phone string) bool
}

type UserStorage struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		pool: pool,
	}
}

func (u *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	userLine := u.pool.QueryRow(ctx, GetUserByEmail, email)

	user := models.User{
		Email: email,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password, &user.Birthday); err != nil {
		log.Printf("error in GetUserByEmail: %+v", err)

		return nil, err
	}

	return &user, nil
}

func (u *UserStorage) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	userLine := u.pool.QueryRow(ctx, GetUserByID, id)

	user := models.User{
		ID: id,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password, &user.Birthday); err != nil {
		log.Printf("error in GetUserByID: %+v", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *UserStorage) IsUserExistByEmail(ctx context.Context, email string, phone string) (bool, error) {
	userRow := u.pool.QueryRow(ctx, IsUserExist, email, phone)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		log.Printf("error in IsEmailBusy: %+v", err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (u *UserStorage) CreateUser(ctx context.Context, preUser *models.UserWithoutID) error {
	_, err := u.pool.Exec(ctx, CreateUser, preUser.Email, preUser.Phone, preUser.Name, preUser.Password, preUser.Birthday)
	if err != nil {
		log.Printf("preUser=%+v Error in CreateUser: %+v", preUser, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

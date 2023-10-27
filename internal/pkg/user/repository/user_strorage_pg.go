package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	GetUserByEmail = `SELECT id, email, phone, name, pass, birthday FROM public."user" WHERE email=$1`
	GetUserByID    = `SELECT id, email, phone, name, pass, birthday FROM public."user" WHERE id=$1`
	IsUserExist    = `SELECT id FROM public."user" WHERE email=$1 OR phone=$2`
	SQLAddUser     = `INSERT INTO public."user" (email, phone, name, pass, birthday) VALUES ($1, $2, $3, $4, $5)`
	SQLGetIDUser   = `SELECT id FROM public."user" WHERE email=$1`
)

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

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Pass, &user.Birthday); err != nil {
		log.Printf("error in GetUserByEmail: %v", err)

		return nil, err
	}

	return &user, nil
}

func (u *UserStorage) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	userLine := u.pool.QueryRow(ctx, GetUserByID, id)
	user := models.User{
		ID: id,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Pass, &user.Birthday); err != nil {
		log.Printf("error in GetUserByID: %v", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *UserStorage) IsUserExist(ctx context.Context, email string, phone string) (bool, error) {
	userRow := u.pool.QueryRow(ctx, IsUserExist, email, phone)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		log.Printf("error in IsUserExist: %v", err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (u *UserStorage) AddUser(ctx context.Context, preUser *models.UserWithoutID) (*models.User, error) {
	user := models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		_, err := u.pool.Exec(ctx, SQLAddUser, preUser.Email, preUser.Phone,
			preUser.Name, preUser.Pass, preUser.Birthday)
		if err != nil {
			log.Printf("preUser=%+v Error in AddUser: %v", preUser, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		row := u.pool.QueryRow(ctx, SQLGetIDUser, preUser.Email)

		err = row.Scan(&user.ID)
		if err != nil {
			log.Printf("error in AddUser: %v", err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Email = preUser.Email
	user.Phone = preUser.Phone
	user.Name = preUser.Name
	user.Pass = preUser.Pass
	user.Birthday = preUser.Birthday

	return &user, nil
}

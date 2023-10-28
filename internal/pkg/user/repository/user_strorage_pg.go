package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	CreateUser(ctx context.Context, preUser *models.UserWithoutID) error
	UpdateUser(ctx context.Context, userID uint64, updateData map[string]interface{}) error
	IsEmailBusy(ctx context.Context, email string) (bool, error)
	IsPhoneBusy(ctx context.Context, phone string) (bool, error)
}

type UserStorage struct {
	Pool *pgxpool.Pool
}

func NewUserStorage(Pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		Pool: Pool,
	}
}

func (u *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	GetUserByEmail := `SELECT id, email, phone, name, pass, birthday FROM public."user" WHERE email=$1`
	userLine := u.Pool.QueryRow(ctx, GetUserByEmail, email)

	user := models.User{
		Email: email,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Pass, &user.Birthday); err != nil {
		log.Printf("error in GetUserByEmail: %+v", err)

		return nil, err
	}

	return &user, nil
}

func (u *UserStorage) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	GetUserByID := `SELECT id, email, phone, name, pass, birthday FROM public."user" WHERE id=$1`
	userLine := u.Pool.QueryRow(ctx, GetUserByID, id)

	user := models.User{
		ID: id,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Pass, &user.Birthday); err != nil {
		log.Printf("error in GetUserByID: %+v", err)

		return nil, fmt.Errorf("%w", err)
	}

	return &user, nil
}

func (u *UserStorage) IsEmailBusy(ctx context.Context, email string) (bool, error) {
	IsEmailBusy := `SELECT id FORM public."user" WHERE email=$1`
	userRow := u.Pool.QueryRow(ctx, IsEmailBusy, email)

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

func (u *UserStorage) IsPhoneBusy(ctx context.Context, phone string) (bool, error) {
	IsPhoneBusy := `SELECT id FORM public."user" WHERE phone=$1`
	userRow := u.Pool.QueryRow(ctx, IsPhoneBusy, phone)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		log.Printf("error in IsPhoneBusy: %+v", err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (u *UserStorage) CreateUser(ctx context.Context, preUser *models.UserWithoutID) error {
	CreateUser := `INSERT INTO public."user" (email, phone, name, pass, birthday) VALUES ($1, $2, $3, $4, $5)`
	_, err := u.Pool.Exec(ctx, CreateUser, preUser.Email, preUser.Name, preUser.Name, preUser.Pass, preUser.Phone)
	if err != nil {
		log.Printf("preUser=%+v Error in CreateUser: %+v", preUser, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (u *UserStorage) UpdateUser(ctx context.Context, userID uint64, updateData map[string]interface{}) error {
	query := squirrel.Update(`public."user"`).
		Where(squirrel.Eq{"id": userID})

	for key, value := range updateData {
		query = query.Set(key, value)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error in UpdateUser while converting ToSql: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	_, err = u.Pool.Exec(ctx, queryString, args...)
	if err != nil {
		log.Printf("Error in UpdateUser while executing: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (u *UserStorage) AddUser(ctx context.Context, preUser *models.UserWithoutID) (*models.User, error) {
	user := models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.Pool, func(tx pgx.Tx) error {
		_, err := u.Pool.Exec(ctx, SQLAddUser, preUser.Email, preUser.Phone,
			preUser.Name, preUser.Pass, preUser.Birthday)
		if err != nil {
			log.Printf("preUser=%+v Error in AddUser: %v", preUser, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		row := u.Pool.QueryRow(ctx, SQLGetIDUser, preUser.Email)

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

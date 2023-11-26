package repository

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AuthStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewAuthStorage(pool *pgxpool.Pool) (*AuthStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &AuthStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (a *AuthStorage) isEmailBusy(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	SQLIsEmailBusy := `SELECT id FROM public."user" WHERE email=$1;`
	userRow := tx.QueryRow(ctx, SQLIsEmailBusy, email)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		a.logger.Errorln(err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (a *AuthStorage) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	SQLGetUserByEmail := `SELECT id, email, password FROM public."user" WHERE email=$1;`
	userLine := tx.QueryRow(ctx, SQLGetUserByEmail, email)

	user := models.User{ //nolint:exhaustruct
		Email: email,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password, &user.Birthday); err != nil {
		a.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *AuthStorage) GetUser(ctx context.Context, email string, password string) (*models.User, error) {
	user := &models.User{}            //nolint:exhaustruct
	userWithoutPass := &models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		emailBusy, err := u.isEmailBusy(ctx, tx, email)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !emailBusy {
			return ErrEmailNotExist
		}

		user, err = u.getUserByEmail(ctx, tx, email)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		hashPass, err := hex.DecodeString(user.Password)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !utils.ComparePassAndHash(hashPass, password) {
			return ErrWrongCredentials
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPass.ID = user.ID
	userWithoutPass.Email = user.Email
	userWithoutPass.Phone = user.Phone
	userWithoutPass.Name = user.Name
	userWithoutPass.Birthday = user.Birthday

	return userWithoutPass, nil
}

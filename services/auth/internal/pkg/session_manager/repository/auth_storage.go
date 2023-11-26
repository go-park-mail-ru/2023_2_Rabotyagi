package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrEmailNotExist = myerrors.NewErrorBadContentRequest("Такой email не существует")
	ErrEmailBusy     = myerrors.NewErrorBadContentRequest("Такой email уже занят")

	NameSeqUser = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
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

func GetLastValSeq(ctx context.Context, tx pgx.Tx, logger *zap.SugaredLogger, nameTable pgx.Identifier) (uint64, error) {
	sanitizedNameTable := nameTable.Sanitize()
	SQLGetLastValSeq := fmt.Sprintf(`SELECT last_value FROM %s;`, sanitizedNameTable)
	seqRow := tx.QueryRow(ctx, SQLGetLastValSeq)

	var count uint64

	if err := seqRow.Scan(&count); err != nil {
		logger.Errorf("error in GetLastValSeq: %+v", err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return count, nil
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

	if err := userLine.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		a.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *AuthStorage) GetUser(ctx context.Context, email string) (*models.User, error) {
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

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPass.ID = user.ID
	userWithoutPass.Email = user.Email

	return userWithoutPass, nil
}

func (u *AuthStorage) createUser(ctx context.Context, tx pgx.Tx, email string, password string) error {
	var SQLCreateUser string

	var err error

	SQLCreateUser = `INSERT INTO public."user" (email, password) VALUES ($1, $2);`
	_, err = tx.Exec(ctx, SQLCreateUser, email, password)

	if err != nil {
		u.logger.Errorf("in createUser: preUser=%+v err=%+v", email, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (u *AuthStorage) AddUser(ctx context.Context, email string, password string) (*models.User, error) {
	user := models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		emailBusy, err := u.isEmailBusy(ctx, tx, email)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if emailBusy {
			return ErrEmailBusy
		}

		err = u.createUser(ctx, tx, email, password)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		id, err := GetLastValSeq(ctx, tx, u.logger, NameSeqUser)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		user.ID = id

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Email = email
	user.Password = password

	return &user, nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/pgxpool"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrEmailNotExist = myerrors.NewErrorBadContentRequest("Такой email не существует")
	ErrEmailBusy     = myerrors.NewErrorBadContentRequest("Такой email уже занят")

	NameSeqUser = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
)

type AuthStorage struct {
	pool   pgxpool.IPgxPool
	logger *mylogger.MyLogger
}

func NewAuthStorage(pool pgxpool.IPgxPool) (*AuthStorage, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &AuthStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func GetLastValSeq(ctx context.Context,
	tx pgx.Tx, logger *mylogger.MyLogger, nameTable pgx.Identifier,
) (uint64, error) {
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
	logger := a.logger.LogReqID(ctx)

	SQLIsEmailBusy := `SELECT id FROM public."user" WHERE email=$1;`
	userRow := tx.QueryRow(ctx, SQLIsEmailBusy, email)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		logger.Errorln(err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (a *AuthStorage) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	logger := a.logger.LogReqID(ctx)

	SQLGetUserByEmail := `SELECT id, email, password FROM public."user" WHERE email=$1;`
	userLine := tx.QueryRow(ctx, SQLGetUserByEmail, email)

	user := models.User{ //nolint:exhaustruct
		Email: email,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (a *AuthStorage) GetUser(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	err := pgx.BeginFunc(ctx, a.pool, func(tx pgx.Tx) error {
		emailBusy, err := a.isEmailBusy(ctx, tx, email)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !emailBusy {
			return ErrEmailNotExist
		}

		user, err = a.getUserByEmail(ctx, tx, email)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return user, nil
}

func (a *AuthStorage) createUser(ctx context.Context, tx pgx.Tx, email string, password string) error {
	var SQLCreateUser string

	var err error

	logger := a.logger.LogReqID(ctx)

	SQLCreateUser = `INSERT INTO public."user" (email, password) VALUES ($1, $2);`
	_, err = tx.Exec(ctx, SQLCreateUser, email, password)

	if err != nil {
		logger.Errorf("in createUser: preUser=%+v err=%+v", email, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (a *AuthStorage) AddUser(ctx context.Context, email string, password string) (*models.User, error) {
	user := models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, a.pool, func(tx pgx.Tx) error {
		emailBusy, err := a.isEmailBusy(ctx, tx, email)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if emailBusy {
			return ErrEmailBusy
		}

		err = a.createUser(ctx, tx, email, password)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		id, err := GetLastValSeq(ctx, tx, a.logger, NameSeqUser)
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

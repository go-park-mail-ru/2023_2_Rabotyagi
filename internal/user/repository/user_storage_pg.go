package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/pgxpool"
	"github.com/jackc/pgx/v5"
)

var (
	ErrEmailBusy          = myerrors.NewErrorBadContentRequest("Такой email уже занят")
	ErrEmailNotExist      = myerrors.NewErrorBadContentRequest("Такой email не существует")
	ErrPhoneBusy          = myerrors.NewErrorBadContentRequest("Такой телефон уже занят")
	ErrWrongCredentials   = myerrors.NewErrorBadContentRequest("Некорректный логин или пароль")
	ErrNoUpdateFields     = myerrors.NewErrorBadFormatRequest("Вы пытаетесь обновить пустое количество полей")
	ErrNoAffectedUserRows = myerrors.NewErrorBadFormatRequest("Не получилось обновить данные пользователя")

	NameSeqUser = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
)

type UserStorage struct {
	pool   pgxpool.IPgxPool
	logger *mylogger.MyLogger
}

func NewUserStorage(pool pgxpool.IPgxPool) (*UserStorage, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &UserStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (u *UserStorage) getUserWithoutPasswordByID(ctx context.Context,
	tx pgx.Tx, userID uint64,
) (*models.UserWithoutPassword, error) {
	logger := u.logger.LogReqID(ctx)

	SQLGetUserByID := `SELECT email, phone, name, birthday, avatar, created_at FROM public."user" WHERE id=$1;`
	userLine := tx.QueryRow(ctx, SQLGetUserByID, userID)
	user := models.UserWithoutPassword{ //nolint:exhaustruct
		ID: userID,
	}

	if err := userLine.Scan(&user.Email,
		&user.Phone, &user.Name, &user.Birthday, &user.Avatar, &user.CreatedAt); err != nil {
		logger.Errorf("error in getUserWithoutPasswordByID: %+v", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *UserStorage) getAvgRatingUserByID(ctx context.Context,
	tx pgx.Tx, userID uint64,
) (sql.NullFloat64, error) {
	logger := u.logger.LogReqID(ctx)

	SQLGetAvgRatingUserByID := `SELECT AVG(rating)
FROM public."comment"
WHERE recipient_id = $1;`
	avgRatingLine := tx.QueryRow(ctx, SQLGetAvgRatingUserByID, userID)

	var avgRating sql.NullFloat64

	if err := avgRatingLine.Scan(&avgRating); err != nil {
		logger.Errorf("error in getAvgRatingUserByID: %+v", err)

		return sql.NullFloat64{}, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return avgRating, nil
}

func (u *UserStorage) GetUserWithoutPasswordByID(ctx context.Context, id uint64) (*models.UserWithoutPassword, error) { //nolint:varnamelen
	var user *models.UserWithoutPassword

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		userInner, err := u.getUserWithoutPasswordByID(ctx, tx, id)
		user = userInner
		if err != nil {
			return err
		}

		avgRating, err := u.getAvgRatingUserByID(ctx, tx, id)
		user.AvgRating = avgRating

		return err
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return user, nil
}

func (u *UserStorage) updateUser(ctx context.Context,
	tx pgx.Tx, userID uint64, updateDataMap map[string]interface{},
) error {
	logger := u.logger.LogReqID(ctx)

	if len(updateDataMap) == 0 {
		return ErrNoUpdateFields
	}

	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(`public."user"`).
		Where(squirrel.Eq{"id": userID}).SetMap(updateDataMap)

	queryString, args, err := query.ToSql()
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	result, err := tx.Exec(ctx, queryString, args...)
	if err != nil {
		logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedUserRows)
	}

	return nil
}

func (u *UserStorage) UpdateUser(ctx context.Context,
	userID uint64, updateData map[string]interface{},
) (*models.UserWithoutPassword, error) {
	userWithoutPass := &models.UserWithoutPassword{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		err := u.updateUser(ctx, tx, userID, updateData)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		userWithoutPass, err = u.getUserWithoutPasswordByID(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return userWithoutPass, nil
}

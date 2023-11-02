package repository

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailBusy     = myerrors.NewError("same email already in use")
	ErrPhoneBusy     = myerrors.NewError("same phone already in use")
	ErrWrongPassword = myerrors.NewError("password is wrong")

	NameSeqUser = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
)

type UserStorage struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		pool: pool,
	}
}

func (u *UserStorage) getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	SQLGetUserByEmail := `SELECT id, email, phone, name, password, birthday FROM public."user" WHERE email=$1;`
	userLine := tx.QueryRow(ctx, SQLGetUserByEmail, email)

	user := models.User{ //nolint:exhaustruct
		Email: email,
	}

	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password, &user.Birthday); err != nil {
		log.Printf("error in getUserByEmail: %+v", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		userInner, err := u.getUserByEmail(ctx, tx, email)
		user = userInner

		return err
	})

	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return user, nil
}

func (u *UserStorage) getUserWithoutPasswordByID(ctx context.Context, tx pgx.Tx, id uint64) (*models.UserWithoutPassword, error) { //nolint:lll
	SQLGetUserByID := `SELECT email, phone, name, birthday FROM public."user" WHERE id=$1;`
	userLine := tx.QueryRow(ctx, SQLGetUserByID, id)

	user := models.UserWithoutPassword{ //nolint:exhaustruct
		ID: id,
	}

	if err := userLine.Scan(&user.Email, &user.Phone, &user.Name, &user.Birthday); err != nil {
		log.Printf("error in GetUserWithoutPasswordByID: %+v", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *UserStorage) GetUserWithoutPasswordByID(ctx context.Context, id uint64) (*models.UserWithoutPassword, error) {
	var user *models.UserWithoutPassword

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		userInner, err := u.getUserWithoutPasswordByID(ctx, tx, id)
		user = userInner

		return err
	})

	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return user, nil
}

func (u *UserStorage) isEmailBusy(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	SQLIsEmailBusy := `SELECT id FROM public."user" WHERE email=$1;`
	userRow := tx.QueryRow(ctx, SQLIsEmailBusy, email)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		log.Printf("error in isEmailBusy: %+v", err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (u *UserStorage) IsEmailBusy(ctx context.Context, email string) (bool, error) {
	var emailBusy bool

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		emailBusyInner, err := u.isEmailBusy(ctx, tx, email)
		emailBusy = emailBusyInner

		return err
	})

	return emailBusy, fmt.Errorf(myerrors.ErrTemplate, err)
}

func (u *UserStorage) isPhoneBusy(ctx context.Context, tx pgx.Tx, phone string) (bool, error) {
	SQLIsPhoneBusy := `SELECT id FROM public."user" WHERE phone=$1;`
	userRow := tx.QueryRow(ctx, SQLIsPhoneBusy, phone)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		log.Printf("error in isPhoneBusy: %+v", err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (u *UserStorage) IsPhoneBusy(ctx context.Context, phone string) (bool, error) {
	var phoneBusy bool

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		phoneBusyInner, err := u.isPhoneBusy(ctx, tx, phone)
		phoneBusy = phoneBusyInner

		return err
	})

	return phoneBusy, fmt.Errorf(myerrors.ErrTemplate, err)
}

func (u *UserStorage) createUser(ctx context.Context, tx pgx.Tx, preUser *models.UserWithoutID) error {
	var SQLCreateUser string

	var err error
	preUser.Password, err = utils.HashPass(preUser.Password)

	if err != nil {
		log.Printf("preUser=%+v Error in hashingUser: %+v", preUser, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if preUser.Birthday.IsZero() {
		SQLCreateUser = `INSERT INTO public."user" (email, phone, name, password) VALUES ($1, $2, $3, $4);`
		_, err = tx.Exec(ctx, SQLCreateUser,
			preUser.Email, preUser.Phone, preUser.Name, preUser.Password)
	} else {
		SQLCreateUser = `INSERT INTO public."user" (email, phone, name, password, birthday) VALUES ($1, $2, $3, $4, $5);`
		_, err = tx.Exec(ctx, SQLCreateUser,
			preUser.Email, preUser.Phone, preUser.Name, preUser.Password, preUser.Birthday)
	}

	if err != nil {
		log.Printf("preUser=%+v Error in createUser: %+v", preUser, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (u *UserStorage) updateUser(ctx context.Context,
	tx pgx.Tx, userID uint64, updateDataMap map[string]interface{},
) error {
	query := squirrel.Update(`public."user"`).
		Where(squirrel.Eq{"id": userID})
	query.SetMap(updateDataMap)

	queryString, args, err := query.ToSql()
	if err != nil {
		log.Printf("in updateUser: while converting ToSql: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	_, err = tx.Exec(ctx, queryString, args...)
	if err != nil {
		log.Printf("Error in UpdateUser while executing: %+v", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (u *UserStorage) UpdateUser(ctx context.Context, userID uint64, updateData map[string]interface{}) error {
	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		err := u.updateUser(ctx, tx, userID, updateData)

		return err
	})

	return fmt.Errorf(myerrors.ErrTemplate, err)
}

func (u *UserStorage) getLastValSeq(ctx context.Context, tx pgx.Tx, nameTable pgx.Identifier) (uint64, error) {
	sanitizedNameTable := nameTable.Sanitize()
	SQLGetLastValSeq := fmt.Sprintf(`SELECT last_value FROM %s;`, sanitizedNameTable)
	seqRow := tx.QueryRow(ctx, SQLGetLastValSeq)

	var count uint64

	if err := seqRow.Scan(&count); err != nil {
		log.Printf("error in getLastValSeq: %+v", err)

		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return count, nil
}

func (u *UserStorage) AddUser(ctx context.Context, preUser *models.UserWithoutID) (*models.User, error) {
	user := models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		emailBusy, err := u.isEmailBusy(ctx, tx, preUser.Email)
		if err != nil {
			log.Printf("in AddUser: %+v\n", err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if emailBusy {
			log.Printf("in AddUser: email=%s already busy", preUser.Email)

			return ErrEmailBusy
		}

		phoneBusy, err := u.isPhoneBusy(ctx, tx, preUser.Phone)
		if err != nil {
			log.Printf("in AddUser: %+v\n", err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if phoneBusy {
			log.Printf("in AddUser: email=%s already busy", preUser.Email)

			return ErrPhoneBusy
		}

		err = u.createUser(ctx, tx, preUser)
		if err != nil {
			log.Printf("preUser=%+v Error in AddUser: %+v", preUser, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		id, err := u.getLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {
			log.Printf("in AddUser: %+v", err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		user.ID = id

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Email = preUser.Email
	user.Phone = preUser.Phone
	user.Name = preUser.Name
	user.Password = preUser.Password
	user.Birthday = preUser.Birthday

	return &user, nil
}

//func (u *UserStorage) getUser(ctx context.Context, tx pgx.Tx, email string, password string) (*models.User, error) {
//	SQLGetUserByEmail := `SELECT id, email, phone, name, password, birthday FROM public."user" WHERE email=$1 and password=$2;`
//	userLine := tx.QueryRow(ctx, SQLGetUserByEmail, email, password)
//
//	user := models.User{ //nolint:exhaustruct
//		Email:    email,
//		Password: password,
//	}
//
//	if err := userLine.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password, &user.Birthday); err != nil {
//		log.Printf("error in GetUser: %+v", err)
//
//		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
//	}
//
//	return &user, nil
//}

func (u *UserStorage) GetUser(ctx context.Context, email string, password string) (*models.UserWithoutPassword, error) {
	user := &models.User{}
	userWithoutPass := &models.UserWithoutPassword{}

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		emailBusy, err := u.isEmailBusy(ctx, tx, email)
		if err != nil {
			log.Printf("preUser=%+v Error in GetUser: %+v", email, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !emailBusy {
			log.Printf("in GetUser: email=%s is not exist", email)

			return ErrEmailBusy
		}

		user, err = u.getUserByEmail(ctx, tx, email)
		if err != nil {
			log.Printf("preUser=%+v Error in GetUser: %+v", email, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		hashPass, err := hex.DecodeString(user.Password)
		if err != nil {
			log.Printf("preUser=%+v Error in converting password: %+v", email, err)

			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !utils.ComparePassAndHash(hashPass, password) {
			return ErrWrongPassword
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

package usecases

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"

	"github.com/asaskevich/govalidator"
)

var (
	ErrWrongCredentials = myerrors.NewError("Некорректный логин или пароль")
	ErrDecodeUser       = myerrors.NewError("Некорректный json пользователя")
)

func validateUserWithoutID(r io.Reader) (*models.UserWithoutID, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	decoder := json.NewDecoder(r)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	userWithoutID.Trim()

	_, err = govalidator.ValidateStruct(userWithoutID)

	return userWithoutID, err //nolint:wrapcheck
}

func ValidateUserWithoutID(r io.Reader) (*models.UserWithoutID, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutID, err := validateUserWithoutID(r)
	if err != nil {
		logger.Errorln(err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutID, nil
}

func ValidateUserCredentials(email string, password string) (*models.UserWithoutID, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutID := new(models.UserWithoutID)

	userWithoutID.Email = email
	userWithoutID.Password = password
	userWithoutID.Trim()
	logger.Infoln(userWithoutID)

	_, err = govalidator.ValidateStruct(userWithoutID)
	if err != nil && (govalidator.ErrorByField(err, "email") != "" ||
		govalidator.ErrorByField(err, "password") != "") {
		logger.Errorln(err)

		return nil, ErrWrongCredentials
	}

	return userWithoutID, nil
}

func validateUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	decoder := json.NewDecoder(r)

	userWithoutPassword := new(models.UserWithoutPassword)
	if err := decoder.Decode(userWithoutPassword); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	userWithoutPassword.Trim()

	_, err = govalidator.ValidateStruct(userWithoutPassword)

	return userWithoutPassword, err //nolint:wrapcheck
}

func ValidateUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPassword, err := validateUserWithoutPassword(r)
	if err != nil {
		logger.Errorln(err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutPassword, nil
}

func ValidatePartOfUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPassword, err := validateUserWithoutPassword(r)
	if userWithoutPassword == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		validationErrors := govalidator.ErrorsByField(err)

		for field, err := range validationErrors {
			if err != "non zero value required" {
				logger.Errorln(err)

				return nil, myerrors.NewError("%s error: %s", field, err)
			}
		}
	}

	return userWithoutPassword, nil
}

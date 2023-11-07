package usecases

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

var (
	ErrWrongCredentials = myerrors.NewError("Некорректный логин или пароль")
	ErrDecodeUser       = myerrors.NewError("Некорректный json пользователя")
)

func validateUserWithoutID(logger *zap.SugaredLogger, r io.Reader) (*models.UserWithoutID, error) {
	decoder := json.NewDecoder(r)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		logger.Errorf("in ValidateUserWithoutID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	userWithoutID.Trim()

	_, err := govalidator.ValidateStruct(userWithoutID)

	return userWithoutID, err //nolint:wrapcheck
}

func ValidateUserWithoutID(logger *zap.SugaredLogger, r io.Reader) (*models.UserWithoutID, error) {
	userWithoutID, err := validateUserWithoutID(logger, r)
	if err != nil {
		logger.Errorf("in ValidateUserWithoutID: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutID, nil
}

func ValidateUserCredentials(logger *zap.SugaredLogger, r io.Reader) (*models.UserWithoutID, error) {
	userWithoutID, err := validateUserWithoutID(logger, r)
	if userWithoutID == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil && (govalidator.ErrorByField(err, "email") != "" ||
		govalidator.ErrorByField(err, "password") != "") {
		logger.Errorf("in ValidateUserCredentials: %+v\n", err)

		return nil, ErrWrongCredentials
	}

	return userWithoutID, nil
}

func validateUserWithoutPassword(logger *zap.SugaredLogger, r io.Reader) (*models.UserWithoutPassword, error) {
	decoder := json.NewDecoder(r)

	userWithoutPassword := new(models.UserWithoutPassword)
	if err := decoder.Decode(userWithoutPassword); err != nil {
		logger.Errorf("in ValidateUserWithoutPassword: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	userWithoutPassword.Trim()

	_, err := govalidator.ValidateStruct(userWithoutPassword)

	return userWithoutPassword, err //nolint:wrapcheck
}

func ValidateUserWithoutPassword(logger *zap.SugaredLogger, r io.Reader) (*models.UserWithoutPassword, error) {
	userWithoutPassword, err := validateUserWithoutPassword(logger, r)
	if err != nil {
		logger.Errorf("in ValidateUserWithoutPassword: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutPassword, nil
}

func ValidatePartOfUserWithoutPassword(logger *zap.SugaredLogger, r io.Reader) (*models.UserWithoutPassword, error) {
	userWithoutPassword, err := validateUserWithoutPassword(logger, r)
	if userWithoutPassword == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		validationErrors := govalidator.ErrorsByField(err)

		for field, err := range validationErrors {
			if err != "non zero value required" {
				logger.Errorf("in ValidateUserWithoutPassword: %+v\n", err)

				return nil, myerrors.NewError("%s error: %s", field, err)
			}
		}
	}

	return userWithoutPassword, nil
}

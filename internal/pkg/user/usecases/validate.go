package usecases

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"io"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
)

var (
	ErrDecodeUser = myerrors.NewErrorBadFormatRequest("Некорректный json пользователя")
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

		return nil, myerrors.NewErrorBadFormatRequest(err.Error())
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

	if errMessage := govalidator.ErrorByField(err, "email"); errMessage != "" {
		logger.Errorln(err)

		return nil, myerrors.NewErrorBadContentRequest(errMessage)
	}

	if errMessage := govalidator.ErrorByField(err, "password"); errMessage != "" {
		logger.Errorln(err)

		return nil, myerrors.NewErrorBadContentRequest(errMessage)
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

		return nil, myerrors.NewErrorBadFormatRequest(err.Error())
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

				return nil, myerrors.NewErrorBadFormatRequest("%s error: %s", field, err)
			}
		}
	}

	return userWithoutPassword, nil
}

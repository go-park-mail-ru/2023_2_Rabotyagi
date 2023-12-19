package usecases

import (
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"

	"github.com/asaskevich/govalidator"
)

var (
	ErrDecodeUser                  = myerrors.NewErrorBadFormatRequest("Некорректный json пользователя")
	ErrValidateUserWithoutPassword = myerrors.NewErrorBadContentRequest("Неправильно заполнены поля: ")
)

func validateUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPassword := new(models.UserWithoutPassword)

	data, err := io.ReadAll(r)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	if err := userWithoutPassword.UnmarshalJSON(data); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	userWithoutPassword.Trim()

	_, err = govalidator.ValidateStruct(userWithoutPassword)
	if err != nil {
		logger.Errorln(err)

		// In this place  return non wrapped error because later it should be use in govalidator.ErrorsByField(err)
		return userWithoutPassword, err //nolint:wrapcheck
	}

	return userWithoutPassword, err //nolint:wrapcheck
}

func ValidateUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPassword, err := validateUserWithoutPassword(r)
	if userWithoutPassword == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf("%w %v", ErrValidateUserWithoutPassword, err)
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
		errMessage := ""

		for field, err := range validationErrors {
			if err != "non zero value required" {
				errMessage += fmt.Sprintf("%s error: %s\n", field, err)
			}
		}

		if errMessage != "" {
			logger.Errorln(errMessage)

			return nil, fmt.Errorf("%w %v", ErrValidateUserWithoutPassword, err)
		}
	}

	return userWithoutPassword, nil
}

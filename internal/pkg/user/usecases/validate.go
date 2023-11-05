package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/asaskevich/govalidator"
)

func ValidateUserWithoutID(r io.Reader) (*models.UserWithoutID, error) {
	decoder := json.NewDecoder(r)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		log.Printf("in ValidateUserWithoutID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutID.Email = strings.TrimSpace(userWithoutID.Email)
	userWithoutID.Name = strings.TrimSpace(userWithoutID.Name)
	userWithoutID.Phone = strings.TrimSpace(userWithoutID.Phone)

	_, err := govalidator.ValidateStruct(userWithoutID)
	if err != nil {
		log.Printf("in ValidateUserWithoutID: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutID, nil
}

func validateUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	decoder := json.NewDecoder(r)

	userWithoutPassword := new(models.UserWithoutPassword)
	if err := decoder.Decode(userWithoutPassword); err != nil {
		log.Printf("in ValidateUserWithoutPassword: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPassword.Email = strings.TrimSpace(userWithoutPassword.Email)
	userWithoutPassword.Name = strings.TrimSpace(userWithoutPassword.Name)
	userWithoutPassword.Phone = strings.TrimSpace(userWithoutPassword.Phone)

	_, err := govalidator.ValidateStruct(userWithoutPassword)

	return userWithoutPassword, err //nolint:wrapcheck
}

func ValidateUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	userWithoutPassword, err := validateUserWithoutPassword(r)
	if err != nil {
		log.Printf("in ValidateUserWithoutPassword: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutPassword, nil
}

func ValidatePartOfUserWithoutPassword(r io.Reader) (*models.UserWithoutPassword, error) {
	userWithoutPassword, err := validateUserWithoutPassword(r)
	if userWithoutPassword == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		validationErrors := govalidator.ErrorsByField(err)

		for field, err := range validationErrors {
			if err != "non zero value required" {
				log.Printf("in ValidateUserWithoutPassword: %+v\n", err)

				return nil, myerrors.NewError("%s error: %s", field, err)
			}
		}
	}

	return userWithoutPassword, nil
}

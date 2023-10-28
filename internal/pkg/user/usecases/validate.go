package usecases

import (
	"encoding/json"
	"io"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

func ValidateUserWithoutID(r io.Reader) (*models.UserWithoutID, error) {
	decoder := json.NewDecoder(r)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		log.Printf("in ValidateUserWithoutID: %+v\n", err)

		return nil, err //nolint:wrapcheck
	}

	_, err := govalidator.ValidateStruct(userWithoutID)
	if err != nil {
		log.Printf("in ValidateUserWithoutID: %+v\n", err)

		return nil, errors.NewError(err.Error())
	}

	return userWithoutID, nil
}

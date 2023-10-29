package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

func ValidateUserWithoutID(r io.Reader) (*models.UserWithoutID, error) {
	decoder := json.NewDecoder(r)

	userWithoutID := new(models.UserWithoutID)
	if err := decoder.Decode(userWithoutID); err != nil {
		log.Printf("in ValidateUserWithoutID: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	_, err := govalidator.ValidateStruct(userWithoutID)
	if err != nil {
		log.Printf("in ValidateUserWithoutID: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return userWithoutID, nil
}

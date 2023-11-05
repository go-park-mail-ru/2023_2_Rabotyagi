package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"

	"github.com/asaskevich/govalidator"
)

var ErrDecodePreProduct = myerrors.NewError("incorrect json of Product")

func validatePreProduct(r io.Reader) (*models.PreProduct, error) {
	decoder := json.NewDecoder(r)

	preProduct := new(models.PreProduct)
	if err := decoder.Decode(preProduct); err != nil {
		log.Printf("in ValidatePreProduct: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	preProduct.Trim()

	_, err := govalidator.ValidateStruct(preProduct)

	return preProduct, err //nolint:wrapcheck
}

func ValidatePreProduct(r io.Reader) (*models.PreProduct, error) {
	preProduct, err := validatePreProduct(r)
	if err != nil {
		log.Printf("in ValidatePreProduct: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return preProduct, nil
}

func ValidatePartOfPreProduct(r io.Reader) (*models.PreProduct, error) {
	preProduct, err := validatePreProduct(r)
	if preProduct == nil {
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

	return preProduct, nil
}

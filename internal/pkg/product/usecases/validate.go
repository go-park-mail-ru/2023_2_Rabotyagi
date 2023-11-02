package usecases

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"io"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var ErrDecodePreProduct = myerrors.NewError("incorrect json of Product")

func ValidatePreProduct(r io.Reader) (*models.PreProduct, error) {
	decoder := json.NewDecoder(r)

	preProduct := new(models.PreProduct)
	if err := decoder.Decode(preProduct); err != nil {
		log.Printf("in ValidatePreProduct: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	preProduct.Trim()

	_, err := govalidator.ValidateStruct(preProduct)
	if err != nil {
		log.Printf("in ValidatePreProduct: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return preProduct, nil
}

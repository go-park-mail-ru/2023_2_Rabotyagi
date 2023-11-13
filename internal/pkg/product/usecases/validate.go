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
	ErrDecodePreProduct   = myerrors.NewError("Некорректный json объявления")
	ErrDecodePreOrder     = myerrors.NewError("Некорректный json заказа")
	ErrDecodeOrderChanges = myerrors.NewError("Некорректный json изменения заказа")
	ErrNotExistingStatus  = myerrors.NewError("Статус заказа не может быть больше %d", models.OrderStatusClosed)
)

func validatePreProduct(logger *zap.SugaredLogger, r io.Reader) (*models.PreProduct, error) {
	decoder := json.NewDecoder(r)

	preProduct := &models.PreProduct{
		Delivery: false,
		SafeDeal: false,
	}
	if err := decoder.Decode(preProduct); err != nil {
		logger.Errorf("in ValidatePreProduct: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	preProduct.Trim()

	_, err := govalidator.ValidateStruct(preProduct)

	return preProduct, err //nolint:wrapcheck
}

func ValidatePreProduct(logger *zap.SugaredLogger, r io.Reader) (*models.PreProduct, error) {
	preProduct, err := validatePreProduct(logger, r)
	if err != nil {
		logger.Errorf("in ValidatePreProduct: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	// TODO remove hardcode
	var resultImages []models.Image

	for _, image := range preProduct.Images {
		if image.URL == "" {
			continue
		}

		resultImages = append(resultImages, image)
	}

	preProduct.Images = resultImages

	return preProduct, nil
}

func ValidatePartOfPreProduct(logger *zap.SugaredLogger, r io.Reader) (*models.PreProduct, error) {
	preProduct, err := validatePreProduct(logger, r)
	if preProduct == nil {
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

	return preProduct, nil
}

func ValidatePreOrder(logger *zap.SugaredLogger, r io.Reader) (*models.PreOrder, error) {
	preOrder := new(models.PreOrder)
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(preOrder); err != nil {
		logger.Errorf("in ValidatePreOrder: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreOrder)
	}

	_, err := govalidator.ValidateStruct(preOrder)
	if err != nil {
		logger.Errorf("in ValidatePreOrder: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
	}

	return preOrder, nil
}

func validateOrderChanges(logger *zap.SugaredLogger, r io.Reader) (*models.OrderChanges, error) {
	orderChanges := new(models.OrderChanges)
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(orderChanges); err != nil {
		logger.Errorf("in validateOrderChanges: %+v\n", err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeOrderChanges)
	}

	_, err := govalidator.ValidateStruct(orderChanges)

	return orderChanges, err //nolint:wrapcheck
}

func ValidateOrderChangesCount(logger *zap.SugaredLogger, r io.Reader) (*models.OrderChanges, error) {
	orderChanges, err := validateOrderChanges(logger, r)
	if orderChanges == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		errID := govalidator.ErrorByField(err, "id")
		errCount := govalidator.ErrorByField(err, "count")

		if errID != "" || errCount != "" {
			return nil, myerrors.NewError("%s\n%s", errCount, errID)
		}
	}

	return orderChanges, nil
}

func ValidateOrderChangesStatus(logger *zap.SugaredLogger, r io.Reader) (*models.OrderChanges, error) {
	orderChanges, err := validateOrderChanges(logger, r)
	if orderChanges == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		errStatus := govalidator.ErrorByField(err, "status")
		errID := govalidator.ErrorByField(err, "id")

		if errID != "" || errStatus != "" {
			return nil, myerrors.NewError("%s\n%s", errStatus, errID)
		}
	}

	if orderChanges.Status > models.OrderStatusClosed {
		return nil, ErrNotExistingStatus
	}

	return orderChanges, nil
}

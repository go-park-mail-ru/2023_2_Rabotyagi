package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"

	"github.com/asaskevich/govalidator"
)

var (
	ErrDecodePreProduct   = myerrors.NewErrorBadFormatRequest("Некорректный json объявления")
	ErrDecodePreOrder     = myerrors.NewErrorBadFormatRequest("Некорректный json заказа")
	ErrDecodeOrderChanges = myerrors.NewErrorBadFormatRequest("Некорректный json изменения заказа")
	ErrNotExistingStatus  = myerrors.NewErrorBadFormatRequest(
		"Статус заказа не может быть больше %d", models.OrderStatusClosed)
)

func validatePreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r)
	preProduct := &models.PreProduct{ //nolint:exhaustruct
		Delivery: false,
		SafeDeal: false,
	}

	if err := decoder.Decode(preProduct); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	preProduct.Trim()

	preProduct.SalerID = userID

	_, err = govalidator.ValidateStruct(preProduct)
	if err != nil {
		logger.Errorln(err)

		// In this place  return non wrapped error because later it should be use in govalidator.ErrorsByField(err)
		return preProduct, err //nolint:wrapcheck
	}

	return preProduct, nil
}

func ValidatePreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	preProduct, err := validatePreProduct(r, userID)
	if err != nil {
		myErr := &myerrors.Error{}
		if errors.As(err, &myErr) {
			return nil, fmt.Errorf(myerrors.ErrTemplate, err)
		}

		return nil, myerrors.NewErrorBadContentRequest(err.Error())
	}

	return preProduct, nil
}

func ValidatePartOfPreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	preProduct, err := validatePreProduct(r, userID)
	if preProduct == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		validationErrors := govalidator.ErrorsByField(err)

		for field, err := range validationErrors {
			if err != "non zero value required" {
				logger.Errorln(err)

				return nil, myerrors.NewErrorBadContentRequest("в поле %s ошибка: %s", field, err)
			}
		}
	}

	return preProduct, nil
}

func ValidatePreOrder(r io.Reader) (*models.PreOrder, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	preOrder := new(models.PreOrder)
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(preOrder); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreOrder)
	}

	_, err = govalidator.ValidateStruct(preOrder)
	if err != nil {
		logger.Errorln(err)

		return nil, myerrors.NewErrorBadContentRequest(err.Error())
	}

	return preOrder, nil
}

func validateOrderChanges(r io.Reader) (*models.OrderChanges, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	orderChanges := new(models.OrderChanges)
	decoder := json.NewDecoder(r)

	if err = decoder.Decode(orderChanges); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeOrderChanges)
	}

	_, err = govalidator.ValidateStruct(orderChanges)
	if err != nil {
		logger.Errorln(err)

		// In this place  return non wrapped error because later it should be use in govalidator.ErrorsByField(err)
		return orderChanges, err //nolint:wrapcheck
	}

	return orderChanges, nil
}

func ValidateOrderChangesCount(r io.Reader) (*models.OrderChanges, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	orderChanges, err := validateOrderChanges(r)
	if orderChanges == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		errID := govalidator.ErrorByField(err, "id")
		errCount := govalidator.ErrorByField(err, "count")

		if errID != "" || errCount != "" {
			errInner := myerrors.NewErrorBadFormatRequest("%s\n%s", errCount, errID)
			logger.Errorln(errInner)

			return nil, errInner
		}
	}

	return orderChanges, nil
}

func ValidateOrderChangesStatus(r io.Reader) (*models.OrderChanges, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	orderChanges, err := validateOrderChanges(r)
	if orderChanges == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		errStatus := govalidator.ErrorByField(err, "status")
		errID := govalidator.ErrorByField(err, "id")

		if errID != "" || errStatus != "" {
			errInner := myerrors.NewErrorBadFormatRequest("%s\n%s", errStatus, errID)
			logger.Errorln(errInner)

			return nil, errInner
		}
	}

	if orderChanges.Status > models.OrderStatusClosed {
		errInner := fmt.Errorf(myerrors.ErrTemplate, ErrNotExistingStatus)
		logger.Errorln(errInner)

		return nil, errInner
	}

	return orderChanges, nil
}

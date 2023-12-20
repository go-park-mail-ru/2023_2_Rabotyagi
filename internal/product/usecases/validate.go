package usecases

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"io"

	"github.com/asaskevich/govalidator"
)

var (
	ErrDecodePreProduct   = myerrors.NewErrorBadFormatRequest("Некорректный json объявления")
	ErrDecodePreOrder     = myerrors.NewErrorBadFormatRequest("Некорректный json заказа")
	ErrDecodeProductID    = myerrors.NewErrorBadFormatRequest("Некорректный json product_id")
	ErrDecodeOrderChanges = myerrors.NewErrorBadFormatRequest("Некорректный json изменения заказа")
	ErrNotExistingStatus  = myerrors.NewErrorBadFormatRequest(
		"Статус заказа не может быть больше %d", models.OrderStatusClosed)
	ErrValidatePreProduct         = myerrors.NewErrorBadContentRequest("Ошибка валидации объявления: ")
	ErrValidatePreOrder           = myerrors.NewErrorBadContentRequest("Ошибка валидации заказа: ")
	ErrValidateOrderChangesCount  = myerrors.NewErrorBadFormatRequest("Ошибка валидации количества изменения заказа: ")
	ErrValidateOrderChangesStatus = myerrors.NewErrorBadFormatRequest("Ошибка валидации статуса изменения заказа: ")
)

func validatePreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	preProduct := &models.PreProduct{ //nolint:exhaustruct
		Delivery: false,
		SafeDeal: false,
	}

	data, err := io.ReadAll(r)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	if err := preProduct.UnmarshalJSON(data); err != nil {
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

		return nil, fmt.Errorf("%w %w", ErrValidatePreProduct, err)
	}

	return preProduct, nil
}

func ValidatePartOfPreProduct(r io.Reader, userID uint64) (*models.PreProduct, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
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

				return nil, fmt.Errorf("%w в поле %s ошибка: %s",
					ErrValidatePreProduct, field, err)
			}
		}
	}

	return preProduct, nil
}

func ValidatePreOrder(r io.Reader) (*models.PreOrder, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	preOrder := new(models.PreOrder)

	data, err := io.ReadAll(r)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreOrder)
	}

	if err := preOrder.UnmarshalJSON(data); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreOrder)
	}

	_, err = govalidator.ValidateStruct(preOrder)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf("%w %v", ErrValidatePreOrder, err) //nolint:errorlint
	}

	return preOrder, nil
}

func validateOrderChanges(r io.Reader) (*models.OrderChanges, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	orderChanges := new(models.OrderChanges)

	data, err := io.ReadAll(r)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeOrderChanges)
	}

	if err := orderChanges.UnmarshalJSON(data); err != nil {
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
	logger, err := mylogger.Get()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	orderChanges, err := validateOrderChanges(r)
	if orderChanges == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		errID := govalidator.ErrorByField(err, "id")
		errCount := govalidator.ErrorByField(err, "count")

		if errID != "" || errCount != "" {
			errInner := fmt.Errorf("%w %s\n%s", ErrValidateOrderChangesCount, errCount, errID)
			logger.Errorln(errInner)

			return nil, errInner
		}
	}

	return orderChanges, nil
}

func ValidateOrderChangesStatus(r io.Reader) (*models.OrderChanges, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	orderChanges, err := validateOrderChanges(r)
	if orderChanges == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		errStatus := govalidator.ErrorByField(err, "status")
		errID := govalidator.ErrorByField(err, "id")

		if errID != "" || errStatus != "" {
			errInner := fmt.Errorf("%w %s\n%s", ErrValidateOrderChangesStatus, errStatus, errID)
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

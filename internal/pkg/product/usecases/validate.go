package usecases

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

var (
	ErrDecodePreProduct   = myerrors.NewError("Некорректный json объявления")
	ErrDecodePreOrder     = myerrors.NewError("Некорректный json заказа")
	ErrDecodeOrderChanges = myerrors.NewError("Некорректный json изменения заказа")
	ErrNotExistingStatus  = myerrors.NewError("Статус заказа не может быть больше %d", models.OrderStatusClosed)
)

func validatePreProduct(r io.Reader) (*models.PreProduct, error) {
	decoder := json.NewDecoder(r)
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	preProduct := &models.PreProduct{
		Delivery: false,
		SafeDeal: false,
	}
	if err := decoder.Decode(preProduct); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreProduct)
	}

	preProduct.Trim()

	_, err = govalidator.ValidateStruct(preProduct)

	return preProduct, err //nolint:wrapcheck
}

func ValidatePreProduct(r io.Reader) (*models.PreProduct, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	preProduct, err := validatePreProduct(r)
	if err != nil {
		logger.Errorln(err)

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
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	preProduct, err := validatePreProduct(r)
	if preProduct == nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	if err != nil {
		validationErrors := govalidator.ErrorsByField(err)

		for field, err := range validationErrors {
			if err != "non zero value required" {
				logger.Errorln(err)

				return nil, myerrors.NewError("%s error: %s", field, err)
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
		logger.Errorf("in ValidatePreOrder: %+v\n", err)

		return nil, myerrors.NewError(err.Error())
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

	if err := decoder.Decode(orderChanges); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeOrderChanges)
	}

	_, err = govalidator.ValidateStruct(orderChanges)

	return orderChanges, err //nolint:wrapcheck
}

func ValidateOrderChangesCount(r io.Reader) (*models.OrderChanges, error) {
	orderChanges, err := validateOrderChanges(r)
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

func ValidateOrderChangesStatus(r io.Reader) (*models.OrderChanges, error) {
	orderChanges, err := validateOrderChanges(r)
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

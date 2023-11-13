package delivery

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

var (
	ErrWrongSalerID   = myerrors.NewError("Получили некорректный saler_id параметр. Он должен быть целым")
	ErrWrongCount     = myerrors.NewError("Получили некорректный count параметр. Он должен быть целым")
	ErrWrongLastID    = myerrors.NewError("Получили некорректный last_id параметр. Он должен быть целым")
	ErrWrongProductID = myerrors.NewError("Получили некорректный product_id параметр. Он должен быть целым")
)

func (p *ProductHandler) createURLToProductFromID(productID uint64) string {
	return fmt.Sprintf("%s%s:%s/api/v1/product/get/%d", p.schema, p.addrOrigin, p.portServer, productID)
}

func parseCountAndLastIDFromRequest(r *http.Request, logger *zap.SugaredLogger) (uint64, uint64, error) {
	countStr := r.URL.Query().Get("count")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w count=%s", ErrWrongCount, countStr)

		logger.Errorf("in parseCountAndLastIDFromRequest: %+v\n", err)

		return 0, 0, err
	}

	lastIDStr := r.URL.Query().Get("last_id")

	lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w last_id=%s", ErrWrongLastID, lastIDStr)
		logger.Errorf("in parseCountAndLastIDFromRequest: %+v\n", err)

		return 0, 0, err
	}

	return count, lastID, nil
}

func parseSalerIDCountLastIDFromRequest(r *http.Request, logger *zap.SugaredLogger) (uint64, uint64, uint64, error) {
	salerIDStr := r.URL.Query().Get("saler_id")

	salerID, err := strconv.ParseUint(salerIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w saler_id=%s", ErrWrongSalerID, salerIDStr)

		logger.Errorf("in parseSalerIDCountLastIDFromRequest: %+v\n", err)

		return 0, 0, 0, err
	}

	countStr := r.URL.Query().Get("count")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w count=%s", ErrWrongCount, countStr)

		logger.Errorf("in parseCountAndLastIDFromRequest: %+v\n", err)

		return 0, 0, 0, err
	}

	lastIDStr := r.URL.Query().Get("last_id")

	lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w last_id=%s", ErrWrongLastID, lastIDStr)
		logger.Errorf("in parseCountAndLastIDFromRequest: %+v\n", err)

		return 0, 0, 0, err
	}

	return salerID, count, lastID, nil
}

func parseIDFromRequest(r *http.Request, logger *zap.SugaredLogger) (uint64, error) {
	IDStr := r.URL.Query().Get("id")

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w id=%s", ErrWrongProductID, IDStr)

		logger.Errorf("in parseIDFromRequest: %+v\n", err)

		return 0, err
	}

	return ID, err
}

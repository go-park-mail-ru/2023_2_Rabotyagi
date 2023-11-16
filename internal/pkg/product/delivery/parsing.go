package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
)

var (
	ErrWrongSalerID   = myerrors.NewError("Получили некорректный saler_id параметр. Он должен быть целым")
	ErrWrongCount     = myerrors.NewError("Получили некорректный count параметр. Он должен быть целым")
	ErrWrongLastID    = myerrors.NewError("Получили некорректный last_id параметр. Он должен быть целым")
	ErrWrongProductID = myerrors.NewError("Получили некорректный product_id параметр. Он должен быть целым")
)

func parseCountAndLastIDFromRequest(r *http.Request) (uint64, uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	countStr := r.URL.Query().Get("count")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w count=%s", ErrWrongCount, countStr)

		logger.Errorln(err)

		return 0, 0, err
	}

	lastIDStr := r.URL.Query().Get("last_id")

	lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w last_id=%s", ErrWrongLastID, lastIDStr)
		logger.Errorln(err)

		return 0, 0, err
	}

	return count, lastID, nil
}

func parseSalerIDCountLastIDFromRequest(r *http.Request) (uint64, uint64, uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, 0, 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	salerIDStr := r.URL.Query().Get("saler_id")

	salerID, err := strconv.ParseUint(salerIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w saler_id=%s", ErrWrongSalerID, salerIDStr)

		logger.Errorln(err)

		return 0, 0, 0, err
	}

	countStr := r.URL.Query().Get("count")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w count=%s", ErrWrongCount, countStr)

		logger.Errorln(err)

		return 0, 0, 0, err
	}

	lastIDStr := r.URL.Query().Get("last_id")

	lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w last_id=%s", ErrWrongLastID, lastIDStr)
		logger.Errorln(err)

		return 0, 0, 0, err
	}

	return salerID, count, lastID, nil
}

func parseIDFromRequest(r *http.Request) (uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		fmt.Println("in parseIDFromRequest: ", err)

		return 0, err
	}

	IDStr := r.URL.Query().Get("id")

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w id=%s", ErrWrongProductID, IDStr)

		logger.Errorln(err)

		return 0, err
	}

	return ID, err
}

package utils

import (
	"fmt"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_logger"
	"net/http"
	"strconv"
)

var ErrWrongNumberParam = myerrors.NewError("Получили некорректный числовой параметр. Он должен быть целым")

func ParseUint64FromRequest(r *http.Request, paramName string) (uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	numberStr := r.URL.Query().Get(paramName)

	number, err := strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
		err := fmt.Errorf("%w %s=%s", ErrWrongNumberParam, paramName, numberStr)

		logger.Errorln(err)

		return 0, err
	}

	return number, err
}

func ParseStringFromRequest(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}

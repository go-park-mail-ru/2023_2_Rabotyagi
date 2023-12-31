package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

var MessageErrWrongNumberParam = "Получили некорректный числовой параметр. " + //nolint:gochecknoglobals
	"Он должен быть целым"

func ParseUint64FromRequest(r *http.Request, paramName string) (uint64, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	numberStr := r.URL.Query().Get(paramName)

	number, err := strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
		err := myerrors.NewErrorBadFormatRequest("%s %s=%s", MessageErrWrongNumberParam, paramName, numberStr)

		logger.Errorln(err)

		return 0, err
	}

	return number, nil
}

func ParseStringFromRequest(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}

package delivery

import (
	"errors"
	"go.uber.org/zap"
	"net/http"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
)

// HandleErr this function handle err. If err is myerror.
// Error then we built this error and client get it, otherwise it is internal error and client shouldn`t get it.
func HandleErr(w http.ResponseWriter, logger *zap.SugaredLogger, err error) {
	myErr := &myerrors.Error{}
	if errors.As(err, &myErr) {
		SendErrResponse(w, logger, NewErrResponse(StatusErrBadRequest, err.Error()))

		return
	}

	SendErrResponse(w, logger, NewErrResponse(StatusErrInternalServer, ErrInternalServer))
}

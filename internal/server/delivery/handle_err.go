package delivery

import (
	"errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/statuses"
	"net/http"

	"go.uber.org/zap"
)

// HandleErr this function handle err. If err is myerror.
// Error then we built this error and client get it, otherwise it is internal error and client shouldn`t get it.
func HandleErr(w http.ResponseWriter, logger *zap.SugaredLogger, err error) {
	myErr := &myerrors.Error{}
	if errors.As(err, &myErr) && myErr.IsErrorClient() {
		SendResponse(w, logger, NewErrResponse(myErr.Status(), err.Error()))

		return
	}

	SendResponse(w, logger, NewErrResponse(statuses.StatusInternalServer, ErrInternalServer))
}

package delivery

import (
	"errors"
	"log"
	"net/http"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
)

// HandleErr this function handle err. If err is myerror.
// Error then we built this error and client get it, otherwise it is internal error and client shouldn`t get it.
func HandleErr(w http.ResponseWriter, message string, err error) {
	log.Printf("%s %+v\n", message, err)

	myErr := &myerrors.Error{}
	if errors.As(err, &myErr) {
		SendErrResponse(w, NewErrResponse(StatusErrBadRequest, err.Error()))

		return
	}

	SendErrResponse(w, NewErrResponse(StatusErrInternalServer, ErrInternalServer))
}

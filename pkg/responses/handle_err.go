package responses

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"

	"google.golang.org/grpc/status"
)

// HandleErr this function handle err. If err is myerror.
// Error then we built this error and client get it, otherwise it is internal error and client shouldn`t get it.
func HandleErr(w http.ResponseWriter, logger *my_logger.MyLogger, err error) {
	myErr := &myerrors.Error{}
	if errors.As(err, &myErr) && myErr.IsErrorClient() {
		SendResponse(w, logger, NewErrResponse(myErr.Status(), err.Error()))

		return
	}

	if grpcErr := status.Convert(err); grpcErr != nil {
		myErr = myerrors.NewErrorCustom(int(grpcErr.Code()), grpcErr.Message())
		if myErr.IsErrorClient() {
			SendResponse(w, logger, NewErrResponse(myErr.Status(), grpcErr.Message()))

			return
		}
	}

	SendResponse(w, logger, NewErrResponse(statuses.StatusInternalServer, ErrInternalServer))
}

package responses

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"google.golang.org/grpc/status"
)

func trimErrMessage(message string) string {
	endInternal := "desc = "
	idxEndInternal := strings.Index(message, endInternal)

	if len(message) > idxEndInternal+len(endInternal) {
		return message[idxEndInternal+len(endInternal):]
	}

	return message
}

// HandleErr this function handle err. If err is myerror.
// Error then we built this error and client get it, otherwise it is internal error and client shouldn`t get it.
// Also hear added status in ctx of request
func HandleErr(w http.ResponseWriter, request *http.Request, logger *mylogger.MyLogger, err error) {
	myErr := &myerrors.Error{}
	if errors.As(err, &myErr) && myErr.IsErrorClient() {
		*request = *request.WithContext(statuses.FillStatusCtx(request.Context(), myErr.Status()))
		SendResponse(w, logger, NewErrResponse(myErr.Status(), err.Error()))

		return
	}

	if grpcErr := status.Convert(err); grpcErr != nil {
		myErr = myerrors.NewErrorCustom(int(grpcErr.Code()), grpcErr.Message())
		if myErr.IsErrorClient() {
			*request = *request.WithContext(statuses.FillStatusCtx(request.Context(), myErr.Status()))
			SendResponse(w, logger, NewErrResponse(myErr.Status(), trimErrMessage(grpcErr.Message())))

			return
		}
	}

	*request = *request.WithContext(statuses.FillStatusCtx(request.Context(), myErr.Status()))
	SendResponse(w, logger, NewErrResponse(statuses.StatusInternalServer, ErrInternalServer))
}

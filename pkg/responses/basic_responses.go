package responses

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

const (
	ErrInternalServer = "Ошибка на сервере"
)

var ErrCookieNotPresented = myerrors.NewErrorBadFormatRequest("Должна быть выставлена cookie, а её нет")

const (
	CookieAuthName = "access_token"
)

type Marshaller interface {
	MarshalJSON() ([]byte, error)
}

type ResponseBody struct {
	Message string `json:"message"`
}

//easyjson:json
type ResponseSuccessful struct {
	Status int          `json:"status"`
	Body   ResponseBody `json:"body"`
}

func NewResponseSuccessful(message string) *ResponseSuccessful {
	return &ResponseSuccessful{
		Status: statuses.StatusResponseSuccessful,
		Body:   ResponseBody{message},
	}
}

type ResponseBodyID struct {
	ID uint64 `json:"id"`
}

//easyjson:json
type ResponseID struct {
	Status int            `example:"3003" json:"status"`
	Body   ResponseBodyID `json:"body"`
}

func NewResponseIDRedirect(id uint64) *ResponseID {
	return &ResponseID{Status: statuses.StatusRedirectAfterSuccessful, Body: ResponseBodyID{ID: id}}
}

type ResponseBodyError struct {
	Error string `json:"error"`
}

//easyjson:json
type ErrorResponse struct {
	Status int               `json:"status"`
	Body   ResponseBodyError `json:"body"`
}

func NewErrResponse(status int, err string) *ErrorResponse {
	return &ErrorResponse{
		Status: status,
		Body:   ResponseBodyError{Error: err},
	}
}

func sendResponse(w http.ResponseWriter, logger *my_logger.MyLogger, response Marshaller) {
	responseSend, err := response.MarshalJSON()
	if err != nil {
		logger.Errorf("in sendResponse: %+v\n", err)
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)

		return
	}

	_, err = w.Write(responseSend)
	if err != nil {
		logger.Errorf("in sendResponse: %+v\n", err)
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
	}
}

func SendResponse(w http.ResponseWriter, logger *my_logger.MyLogger, response Marshaller) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sendResponse(w, logger, response)
}

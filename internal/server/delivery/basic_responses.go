package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"net/http"

	"go.uber.org/zap"
)

const (
	ErrInternalServer = "Ошибка на сервере"
)

var ErrCookieNotPresented = myerrors.NewErrorBadFormatRequest("Должна быть выставлена cookie, а её нет")

const (
	CookieAuthName = "access_token"
)

type ResponseBody struct {
	Message string `json:"message"`
}

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

type ResponseID struct {
	Status int            `json:"status"`
	Body   ResponseBodyID `json:"body"`
}

func NewResponseIDRedirect(ID uint64) *ResponseID {
	return &ResponseID{Status: statuses.StatusRedirectAfterSuccessful, Body: ResponseBodyID{ID: ID}}
}

type ResponseBodyError struct {
	Error string `json:"error"`
}

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

func sendResponse(w http.ResponseWriter, logger *zap.SugaredLogger, response any) {
	responseSend, err := json.Marshal(response)
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

func SendResponse(w http.ResponseWriter, logger *zap.SugaredLogger, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sendResponse(w, logger, response)
}

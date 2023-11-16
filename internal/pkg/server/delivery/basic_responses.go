package delivery

import (
	"encoding/json"
	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/my_errors"
	"go.uber.org/zap"
	"net/http"
)

const (
	HTTPStatusOk    = 200
	HTTPStatusError = 222

	StatusResponseSuccessful      = 200
	StatusRedirectAfterSuccessful = 303
	StatusErrBadRequest           = 400
	StatusErrInternalServer       = 500
)

const (
	ErrInternalServer = "Ошибка на сервере"
	ErrBadRequest     = "Некорректный запрос"
)

var ErrCookieNotPresented = myerrors.NewError("Должна быть выставлена cookie, а её нет")

const (
	CookieAuthName = "access_token"
)

type ResponseBody struct {
	Message string `json:"message"`
}

type Response struct {
	Status int          `json:"status"`
	Body   ResponseBody `json:"body"`
}

func NewResponse(status int, message string) *Response {
	return &Response{
		Status: status,
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

func NewResponseID(ID uint64) *ResponseID {
	return &ResponseID{Status: StatusRedirectAfterSuccessful, Body: ResponseBodyID{ID: ID}}
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

func SendErrResponse(w http.ResponseWriter, logger *zap.SugaredLogger, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusError)
	sendResponse(w, logger, response)
}

func SendOkResponse(w http.ResponseWriter, logger *zap.SugaredLogger, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusOk)
	sendResponse(w, logger, response)
}

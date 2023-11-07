package delivery

import (
	"encoding/json"
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

type RedirectBody struct {
	RedirectURL string `json:"redirect_url"`
}

type ResponseRedirect struct {
	Status int          `json:"status"`
	Body   RedirectBody `json:"body"`
}

func NewResponseRedirect(redirectURL string) *ResponseRedirect {
	return &ResponseRedirect{Status: StatusRedirectAfterSuccessful, Body: RedirectBody{RedirectURL: redirectURL}}
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

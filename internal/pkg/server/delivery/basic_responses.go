package delivery

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	HTTPStatusOk    = 200
	HTTPStatusError = 222

	StatusResponseSuccessful = 200
	StatusErrBadRequest      = 400
	StatusErrInternalServer  = 500
)

const (
	ErrInternalServer = "Error in server"
	ErrBadRequest     = "Wrong request"
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

func sendResponse(w http.ResponseWriter, response any) {
	responseSend, err := json.Marshal(response)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)

		return
	}

	_, err = w.Write(responseSend)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
	}
}

func SendErrResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusError)
	sendResponse(w, response)
}

func SendOkResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusOk)
	sendResponse(w, response)
}

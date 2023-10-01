package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

type ResponseBodyError struct {
	Error string `json:"error"`
}

type ErrorResponse struct {
	Status int `json:"status"`
	Body   ResponseBodyError
}

type ResponseBody struct {
	Message string `json:"message"`
}

type Response struct {
	Status int `json:"status"`
	Body   ResponseBody
}

type PostResponse struct {
	Status int `json:"status"`
	Body storage.Post
}

type PostsListResponse struct {
	Status int `json:"status"`
	Body []storage.Post
}

const (
	StatusResponseSuccessful = 200
	StatusErrBadRequest      = 400
	StatusErrServerError     = 500
)

var (
	ResponseSuccessfulSignUp = Response{
		Status: StatusResponseSuccessful, 
		Body: ResponseBody{Message: "Successful sign up"},
	}

	ResponseSuccessfulSignIn = Response{
		Status: StatusResponseSuccessful, 
		Body: ResponseBody{Message: "Successful sign in"},
	}
  
	ErrInternalServer   = ErrorResponse{Status: StatusErrServerError, Body: ResponseBodyError{Error: "Error in server"}}
	ErrBadRequest       = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Wrong request"}}
	ErrUserAlreadyExist = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "User with same email already exist"}}
	ErrWrongCredentials = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Uncorrect login or password"}}

	ErrPostNotExist = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Post not exists"}}
	ErrNoSuchCountOfPosts = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "n > posts count"}}
)

func sendErr(w http.ResponseWriter, errResponse ErrorResponse) {
	response, err := json.Marshal(errResponse)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}

	_, err = w.Write(response)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}
}

func sendResponse(w http.ResponseWriter, response any) {
	responseSend, err := json.Marshal(response)
	if err != nil {
	  log.Printf("%v\n", err)
	  http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}
  
	_, err = w.Write(responseSend)
	if err != nil {
	  log.Printf("%v\n", err)
	  http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}
}

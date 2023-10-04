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
	Body   storage.Post
}

type PostsListResponse struct {
	Status int `json:"status"`
	Body   []storage.Post
}

const (
	StatusResponseSuccessful = 200
	StatusErrBadRequest      = 400
	StatusUnauthorized       = 401
	StatusErrServerError     = 500

	CookieAuthName = "session_id"
)

var (
	ResponseSuccessfulSignUp = Response{Status: StatusResponseSuccessful, Body: ResponseBody{Message: "Successful sign up"}}
	ResponseSuccessfulSignIn = Response{Status: StatusResponseSuccessful, Body: ResponseBody{Message: "Successful sign in"}}
	ResponseSuccessfulLogOut = Response{Status: StatusResponseSuccessful, Body: ResponseBody{Message: "Successful log out"}}

	ResponseSuccessfulAddPost = Response{Status: StatusResponseSuccessful, Body: ResponseBody{Message: "Successful add post"}}

	ErrInternalServer   = ErrorResponse{Status: StatusErrServerError, Body: ResponseBodyError{Error: "Error in server"}}
	ErrBadRequest       = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Wrong request"}}
	ErrUserAlreadyExist = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "User with same email already exist"}}
	ErrWrongCredentials = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Uncorrect login or password"}}
	ErrUnauthorized     = ErrorResponse{Status: StatusUnauthorized, Body: ResponseBodyError{Error: "You unauthorized"}}

	ErrPostNotExist       = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Post not exists"}}
	ErrNoSuchCountOfPosts = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "n > posts count"}}
)

func sendResponse(w http.ResponseWriter, response any) {
	responseSend, err := json.Marshal(response)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)

		return
	}

	_, err = w.Write(responseSend)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, ErrInternalServer.Body.Error, http.StatusInternalServerError)
	}
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://84.23.53.28")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

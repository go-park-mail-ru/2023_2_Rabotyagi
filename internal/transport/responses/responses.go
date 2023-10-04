package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

const CookieAuthName = "session_id"

const (
	HTTPStatusOk    = 200
	HTTPStatusError = 222

	StatusResponseSuccessful = 200
	StatusErrBadRequest      = 400
	StatusUnauthorized       = 401
	StatusErrInternalServer  = 500
)

const (
	ResponseSuccessfulSignUp = "Successful sign up"
	ResponseSuccessfulSignIn = "Successful sign in"
	ResponseSuccessfulLogOut = "Successful log out"

	ResponseSuccessfulAddPost = "Successful add post"

	ErrInternalServer   = "Error in server"
	ErrBadRequest       = "Wrong request"
	ErrUserNotExits     = "User with same email not exist"
	ErrUserAlreadyExist = "User with same email already exist"
	ErrWrongCredentials = "Uncorrected login or password"
	ErrUnauthorized     = "You unauthorized"

	ErrPostNotExist       = "Post not exists"
	ErrNoSuchCountOfPosts = "not enough posts in storage"
)

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

type PostResponse struct {
	Status int           `json:"status"`
	Body   *storage.Post `json:"body"`
}

func NewPostResponse(status int, body *storage.Post) *PostResponse {
	return &PostResponse{
		Status: status,
		Body:   body,
	}
}

type PostsListResponse struct {
	Status int             `json:"status"`
	Body   []*storage.Post `json:"body"`
}

func NewPostsListResponse(status int, body []*storage.Post) *PostsListResponse {
	return &PostsListResponse{
		Status: status,
		Body:   body,
	}
}

// sendResponse don`t use this function, instead use SendOkResponse or SendErrResponse
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
	w.WriteHeader(HTTPStatusError)
	sendResponse(w, response)
}

func SendOkResponse(w http.ResponseWriter, response any) {
	w.WriteHeader(HTTPStatusOk)
	sendResponse(w, response)
}

func SetupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://84.23.53.28")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

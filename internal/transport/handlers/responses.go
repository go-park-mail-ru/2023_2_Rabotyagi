package handler

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

	ErrInternalServer   = ErrorResponse{Status: StatusErrServerError, Body: ResponseBodyError{Error: "Error in server"}}
	ErrBadRequest       = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Wrong request"}}
	ErrUserAlreadyExist = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "User with same email already exist"}}
	ErrWrongCredentials = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Uncorrect login or password"}}
	ErrUnauthorized     = ErrorResponse{Status: StatusUnauthorized, Body: ResponseBodyError{Error: "You unauthorized"}}
)

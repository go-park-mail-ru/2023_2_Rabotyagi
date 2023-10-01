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
	StatusErrServerError     = 500
)

var (
	ResponseSuccessfulSignUp = Response{
		Status: StatusResponseSuccessful,
		Body:   ResponseBody{Message: "Successful sign up"},
	}

	ErrInternalServer   = ErrorResponse{Status: StatusErrServerError, Body: ResponseBodyError{Error: "Error in server"}}
	ErrBadRequest       = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "Wrong request"}}
	ErrUserAlreadyExist = ErrorResponse{Status: StatusErrBadRequest, Body: ResponseBodyError{Error: "User with same email already exist"}}
)

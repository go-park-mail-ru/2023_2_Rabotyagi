package delivery

import "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"

const (
	StatusUnauthorized = 401

	ResponseSuccessfulSignUp = "Successful sign up"
	ResponseSuccessfulSignIn = "Successful sign in"
	ResponseSuccessfulLogOut = "Successful log out"

	ErrWrongCredentials = "Uncorrected login or password"
	ErrUnauthorized     = "You unauthorized"
)

type ProfileResponse struct {
	Status int                         `json:"status"`
	Body   *models.UserWithoutPassword `json:"body"`
}

func NewProfileResponse(status int, body *models.UserWithoutPassword) *ProfileResponse {
	return &ProfileResponse{
		Status: status,
		Body:   body,
	}
}

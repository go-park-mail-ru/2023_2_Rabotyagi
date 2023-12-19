package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

const (
	StatusUnauthorized = 401

	ResponseSuccessfulSignUp = "Successful sign up"
	ResponseSuccessfulSignIn = "Successful sign in"
	ResponseSuccessfulLogOut = "Successful log out"

	ErrUnauthorized = "Вы не авторизованны"
)

//easyjson:json
type ProfileResponse struct {
	Status int                         `json:"status"`
	Body   *models.UserWithoutPassword `json:"body"`
}

func NewProfileResponse(body *models.UserWithoutPassword) *ProfileResponse {
	return &ProfileResponse{
		Status: statuses.StatusResponseSuccessful,
		Body:   body,
	}
}

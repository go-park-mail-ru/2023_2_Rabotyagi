package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

type ResponseURLBody struct {
	SlURL []string `json:"urls"` //nolint:tagliatelle
}

type ResponseURLs struct {
	Status int             `json:"status"`
	Body   ResponseURLBody `json:"body"`
}

func NewResponseURLs(slURL []string) *ResponseURLs {
	return &ResponseURLs{Status: statuses.StatusResponseSuccessful, Body: ResponseURLBody{SlURL: slURL}}
}

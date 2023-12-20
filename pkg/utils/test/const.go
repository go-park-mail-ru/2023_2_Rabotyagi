package test

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
)

// AccessToken for read only, because async usage.
const AccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJlbWFpbCI6Iml2bi0xNS0wN0BtYWlsLnJ1IiwiZXhwaXJlIjoxNzAxMjg1MzE4LCJ1c2VySUQiOjExfQ." +
	"jIPlwcF5xGPpgQ5WYp5kFv9Av-yguX2aOYsAgbodDM4"

// Cookie for read only, because async usage.
var Cookie = http.Cookie{ //nolint:gochecknoglobals,exhaustruct
	Name:  responses.CookieAuthName,
	Value: AccessToken, Expires: time.Now().Add(time.Hour),
}

const UserID uint64 = 1

const ProductID uint64 = 1

const CountProduct uint64 = 2

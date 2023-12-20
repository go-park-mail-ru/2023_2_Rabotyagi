package myerrors

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

const (
	ErrTemplate = "%w"
)

// Error - struct error.
//
// Errors with status values from 4000 to 4999 are considered client errors and sent to the front '
// in internal/pkg/server/delivery/handle_err.go function HandleErr.
// This maybe be checked with help IsErrorClient function.
//
// Errors with status values from 5000 to 5999 are sent to the server in the same HandleErr
// always as internal with status = 5000.
type Error struct {
	err    string
	status int
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) IsErrorClient() bool {
	return e.status >= statuses.MinValueClientError && e.status <= statuses.MaxValueClientError
}

// NewErrorCustom prefer this function to other NewError... from this package. Use this function only if you understand why.
func NewErrorCustom(status int, format string, args ...any) *Error {
	return &Error{status: status, err: fmt.Sprintf(format, args...)}
}

// NewErrorBadFormatRequest error with status =
// StatusBadFormatRequest uses when get bad request from frontend and errors with this status need frontend developer.
func NewErrorBadFormatRequest(format string, args ...any) *Error {
	return &Error{err: fmt.Sprintf(format, args...), status: statuses.StatusBadFormatRequest}
}

// NewErrorBadContentRequest error with status =
// StatusBadContentRequest uses when user has entered incorrect data and needs to show him this error.
func NewErrorBadContentRequest(format string, args ...any) *Error {
	return &Error{err: fmt.Sprintf(format, args...), status: statuses.StatusBadContentRequest}
}

// NewErrorInternal error with status =
// StatusInternalServer uses for indicates internal error status in server.
func NewErrorInternal(format string, args ...any) *Error {
	return &Error{err: fmt.Sprintf(format, args...), status: statuses.StatusInternalServer}
}

func (e *Error) Error() string {
	return e.err
}

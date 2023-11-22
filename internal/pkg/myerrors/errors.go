package myerrors

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery/statuses"
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
// always as internal with status = 5000
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

func NewErrorBadFormatRequest(format string, args ...any) *Error {
	return &Error{err: fmt.Sprintf(format, args...), status: statuses.StatusBadFormatRequest}
}

func NewErrorBadContentRequest(format string, args ...any) *Error {
	return &Error{err: fmt.Sprintf(format, args...), status: statuses.StatusBadContentRequest}
}

func NewErrorInternal(format string, args ...any) *Error {
	return &Error{err: fmt.Sprintf(format, args...), status: statuses.StatusInternalServer}
}

func (e *Error) Error() string {
	return e.err
}

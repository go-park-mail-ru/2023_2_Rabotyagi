package errors

import "fmt"

const (
	ErrTemplate = "%w"
)

type Error struct {
	err string
}

func NewError(format string, args ...any) *Error {
	return &Error{fmt.Sprintf(format, args...)}
}

func (e *Error) Error() string {
	return e.err
}

package error

import (
	"errors"
	"fmt"
)

type Error struct {
	HTTPCode  int
	HTTPError error
	Code      int
	LogError  error
}

func (e *Error) Wrapper(err error) *Error {
	e.LogError = fmt.Errorf("%s\n%w", err.Error(), e.LogError)
	return e
}

func NewError(httpCode int, code int, httpError string, logError error) *Error {
	return &Error{
		HTTPCode:  httpCode,
		HTTPError: errors.New(httpError),
		Code:      code,
		LogError:  logError,
	}
}

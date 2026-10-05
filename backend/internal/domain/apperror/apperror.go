package apperror

import (
	"fmt"
)

type AppError struct {
	Code    string // internal error codes (eg: email_not_verified)
	Status  int    // this is the http status
	Message string
	Err     error
}

func (e *AppError) Error() string {
	err := fmt.Sprintf("Error occurred with code: %s, status: %d, message: %s", e.Code, e.Status, e.Message)
	return err
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(base AppError) *AppError {
	return &base
}

func Wrap(base AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Status:  base.Status,
		Message: base.Message,
		Err:     err,
	}
}

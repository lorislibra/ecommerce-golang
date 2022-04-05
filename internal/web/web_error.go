package web

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type WebError struct {
	Status int
	Err    error
}

func (e *WebError) Error() string {
	return e.Err.Error()
}

func NewWebErrMessage(status int, message string) error {
	return &WebError{
		Status: status,
		Err:    errors.New(message),
	}
}

func NewWebError(status int, err error) error {
	return &WebError{
		Status: status,
		Err:    err,
	}
}

var (
	ErrInternal = NewWebErrMessage(fiber.StatusInternalServerError, "internal error")
	ErrNotFound = NewWebErrMessage(fiber.StatusNotFound, "not found")
)

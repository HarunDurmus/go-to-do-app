package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ansel1/merry/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

func buildErrorResponse(err error) ErrorResponse {
	var errorResponse ErrorResponse
	if errors.As(err, &errorResponse) {
		return errorResponse
	}
	if e, ok := err.(validation.Errors); ok {
		return InvalidInput(e)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return NotFound("")
	}
	if errors.Is(err, fiber.ErrMethodNotAllowed) {
		return MethodNotAllowedError("")
	}

	return InternalServerError("")
}

func New(msg string) error {
	return WithStack(fmt.Errorf(msg))
}

func Wrap(err error, msg string) error {
	return merry.Wrap(err, merry.AppendMessage(msg))
}

func WithStack(err error) error {
	return merry.Wrap(err)
}

func getError(err error, errorResponse ErrorResponse) error {
	if errorResponse.nestedErr != nil {
		return errorResponse.nestedErr
	}
	return err
}

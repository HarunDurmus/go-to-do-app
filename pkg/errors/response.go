package errors

import (
	"net/http"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorDetails struct {
	Description string `json:"description"`
}

type ErrorResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Code      string      `json:"code"`
	Details   interface{} `json:"details,omitempty"`
	nestedErr error       `json:"-"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) StatusCode() int {
	return e.Status
}

func (e ErrorResponse) Err(err error) ErrorResponse {
	e.nestedErr = WithStack(err)
	return e
}

func InternalServerError(msg string) ErrorResponse {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

func MethodNotAllowedError(msg string) ErrorResponse {
	if msg == "" {
		msg = "Method Not Allowed."
	}
	return ErrorResponse{
		Status:  http.StatusMethodNotAllowed,
		Message: msg,
	}
}

func NotFound(msg string) ErrorResponse {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}

}

func Unauthorized(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}

}

func Forbidden(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return ErrorResponse{
		Status:  http.StatusForbidden,
		Message: msg,
	}

}

func BadRequest(msg string) ErrorResponse {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}

}

func Conflict(msg string) ErrorResponse {
	if msg == "" {
		msg = "The request could not be completed due to a conflict with the current state of the target resource."
	}
	return ErrorResponse{
		Status:  http.StatusConflict,
		Message: msg,
	}

}

func PreconditionFailed(msg string) ErrorResponse {
	if msg == "" {
		msg = "One or more conditions given preconditions failed."
	}
	return ErrorResponse{
		Status:  http.StatusPreconditionFailed,
		Message: msg,
	}
}

type InvalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func InvalidInput(errs validation.Errors) ErrorResponse {
	var details []InvalidField
	var fields []string
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		details = append(details, InvalidField{
			Field: field,
			Error: errs[field].Error(),
		})
	}

	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "There is some problem with the data you submitted.",
		Details: details,
	}
}

package v1

import (
	"github.com/savsgio/atreugo/v11"
	"gopkg.in/validator.v2"
)

type ResponseError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Response[K any] struct {
	Status int             `json:"status"`
	Data   *K              `json:"data,omitempty"`
	Errors []ResponseError `json:"errors,omitempty"`
}

func Success[K any](status int, data K) *Response[any] {
	var dataInterface any = data
	return &Response[any]{
		Status: status,
		Data:   &dataInterface,
	}
}

func Errors(status int, errors []ResponseError) *Response[any] {
	return &Response[any]{
		Status: status,
		Errors: errors,
	}
}

func Error(status int, message string) *Response[any] {
	return &Response[any]{
		Status: status,
		Errors: []ResponseError{{Message: message}},
	}
}

func (response *Response[K]) Send(ctx *atreugo.RequestCtx) error {
	return ctx.JSONResponse(response, response.Status)
}

func validatorErrorToResponseError(err error) []ResponseError {
	var responseErrors []ResponseError
	if errs, ok := err.(validator.ErrorMap); ok {
		for field, errors := range errs {
			for _, err := range errors {
				responseErrors = append(responseErrors, ResponseError{
					Field:   field,
					Message: err.Error(),
				})
			}
		}
	}

	return responseErrors
}

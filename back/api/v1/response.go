package v1

import "gopkg.in/validator.v2"

type ResponseError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Response[K any] struct {
	Status int             `json:"status"`
	Data   *K              `json:"data,omitempty"`
	Errors []ResponseError `json:"errors,omitempty"`
}

func NewResponse[K any](status int, data K) Response[K] {
	return Response[K]{
		Status: status,
		Data:   &data,
	}
}

func NewErrorResponse(status int, errors []ResponseError) Response[interface{}] {
	return Response[interface{}]{
		Status: status,
		Errors: errors,
	}
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

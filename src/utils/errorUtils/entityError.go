package errorUtils

import (
	"encoding/json"
	"net/http"
)

type EntityError interface {
	Message() string
	Error() string
	Status() int
}

type entityError struct {
	error
	ErrorMessage string `json:"message"`
	ErrorStatus int `json:"status"`
	ErrError string `json:"error"`
}

func (e *entityError) Error() string {
	return e.ErrError
}

func (e *entityError) Message() string {
	return e.ErrorMessage
}

func (e *entityError) Status() int {
	return e.ErrorStatus
}

func NewEntityError(err error) EntityError {
	if err == nil {
		return nil
	}
	return &entityError{
		ErrorMessage: err.Error(),
		ErrError: err.Error(),
		ErrorStatus: http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) EntityError {
	return &entityError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewBadRequestError(message string) EntityError {
	return &entityError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewUnprocessableEntityError(message string) EntityError {
	return &entityError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusUnprocessableEntity,
		ErrError:   "invalid_request",
	}
}

func NewApiErrFromBytes(body []byte) (EntityError, error) {
	var result entityError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func NewInternalServerError(message string) EntityError {
	return &entityError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusInternalServerError,
		ErrError:   "server_error",
	}
}
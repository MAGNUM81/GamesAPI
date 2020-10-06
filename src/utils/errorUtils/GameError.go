package errorUtils

import (
	"encoding/json"
	"net/http"
)

type GameError interface {
	Message() string
	Error() string
	Status() int
}

type gameError struct {
	ErrorMessage string `json:"message"`
	ErrorStatus int `json:"status"`
	ErrError string `json:"error"`
}

func (e *gameError) Error() string {
	return e.ErrError
}

func (e *gameError) Message() string {
	return e.ErrorMessage
}

func (e *gameError) Status() int {
	return e.ErrorStatus
}

func NewNotFoundError(message string) GameError {
	return &gameError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewBadRequestError(message string) GameError {
	return &gameError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewUnprocessableEntityError(message string) GameError {
	return &gameError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusUnprocessableEntity,
		ErrError:   "invalid_request",
	}
}

func NewApiErrFromBytes(body []byte) (GameError, error) {
	var result gameError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func NewInternalServerError(message string) GameError {
	return &gameError{
		ErrorMessage: message,
		ErrorStatus:  http.StatusInternalServerError,
		ErrError:   "server_error",
	}
}
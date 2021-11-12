package api_error

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ApiErr interface {
	Message() string
	StatusCode() int
	Error() string
	Causes() []interface{}
}

type apiErr struct {
	ErrMessage    string        `json:"message"`
	ErrStatusCode int           `json:"statuscode"`
	ErrCauses     []interface{} `json:"causes"`
}

func (e apiErr) Error() string {
	return fmt.Sprintf("Message: %s - Status: %d - Causes: %v",
		e.ErrMessage, e.ErrStatusCode, e.ErrCauses)
}

func (e apiErr) Message() string {
	return e.ErrMessage
}

func (e apiErr) StatusCode() int {
	return e.ErrStatusCode
}

func (e apiErr) Causes() []interface{} {
	return e.ErrCauses
}

func NewError(msg string, code int, causes []interface{}) ApiErr {
	return apiErr{
		ErrMessage:    msg,
		ErrStatusCode: code,
		ErrCauses:     causes,
	}
}

func NewErrorFromBytes(bytes []byte) (ApiErr, error) {
	var restErr apiErr
	if err := json.Unmarshal(bytes, &restErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return restErr, nil
}

func NewBadRequestError(msg string) ApiErr {
	return apiErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) ApiErr {
	return apiErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(msg string) ApiErr {
	return apiErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusUnauthorized,
	}
}

func NewInternalServerError(msg string, err error) ApiErr {
	result := apiErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusInternalServerError,
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func NewProcessingConflictError(msg string) ApiErr {
	return apiErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusConflict,
	}
}

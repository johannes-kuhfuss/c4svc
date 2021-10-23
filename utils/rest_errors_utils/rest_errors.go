package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RestErr interface {
	Message() string
	StatusCode() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	ErrMessage    string        `json:"message"`
	ErrStatusCode int           `json:"statuscode"`
	ErrCauses     []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("Message: %s - Status: %d - Causes: %v",
		e.ErrMessage, e.ErrStatusCode, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) StatusCode() int {
	return e.ErrStatusCode
}

func (e restErr) Causes() []interface{} {
	return e.ErrCauses
}

func NewError(msg string, code int, causes []interface{}) RestErr {
	return restErr{
		ErrMessage:    msg,
		ErrStatusCode: code,
		ErrCauses:     causes,
	}
}

func NewErrorFromBytes(bytes []byte) (RestErr, error) {
	var restErr restErr
	if err := json.Unmarshal(bytes, &restErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return restErr, nil
}

func NewBadRequestError(msg string) RestErr {
	return restErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) RestErr {
	return restErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(msg string) RestErr {
	return restErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusUnauthorized,
	}
}

func NewInternalServerError(msg string, err error) RestErr {
	result := restErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusInternalServerError,
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func NewProcessingConflictError(msg string) RestErr {
	return restErr{
		ErrMessage:    msg,
		ErrStatusCode: http.StatusConflict,
	}
}

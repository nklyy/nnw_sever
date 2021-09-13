package errors

import (
	"fmt"
	"net/http"
	"nnw_s/pkg/codes"
)

type Error struct {
	Code    codes.Code `json:"code"`
	Status  Status     `json:"status"`
	Message string     `json:"message,omitempty"`
}

func (err Error) Error() string {
	repr := fmt.Sprintf("code: %v; status: %v", err.Code, err.Status)
	if err.Message != "" {
		repr += "; message: " + err.Message
	}
	return repr
}

func New(code codes.Code, status Status) error {
	return &Error{Code: code, Status: status}
}

func WithMessage(target error, msg string, args ...interface{}) error {
	err, ok := target.(*Error)
	if !ok {
		return target
	}
	return &Error{
		Code:    err.Code,
		Status:  err.Status,
		Message: fmt.Sprintf(msg, args...),
	}
}

func HTTPCode(target error) int {
	err, ok := target.(*Error)
	if !ok {
		return http.StatusInternalServerError
	}
	return int(err.Code)
}

func NewInternal(msg string) error {
	return &Error{
		Code:    codes.InternalError,
		Status:  statusInternalError,
		Message: msg,
	}
}

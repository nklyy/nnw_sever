package auth

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidRequest   errors.Status = "invalid_request"
	StatusPermissionDenied errors.Status = "permission_denied"
	StatusUnauthorized     errors.Status = "unauthorized"

	StatusInvalidCode errors.Status = "invalid_code"
	StatusInvalidData errors.Status = "invalid_data"
	StatusInvalidJson errors.Status = "invalid_json"
	StatusWrongToken  errors.Status = "wrong_token"
)

var (
	ErrInvalidRequest   = errors.New(codes.BadRequest, StatusInvalidRequest)
	ErrPermissionDenied = errors.New(codes.Forbidden, StatusPermissionDenied)
	ErrUnauthorized     = errors.New(codes.Unauthorized, StatusUnauthorized)
)

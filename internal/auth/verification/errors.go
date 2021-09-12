package verification

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusCodeAlreadyExists errors.Status = "verification_code_already_exists"
	StatusCodeNotFound      errors.Status = "verification_code_not_found"
	StatusInvalidCode       errors.Status = "verification_code_not_valid"
	StatusInvalidEmail      errors.Status = "invalid_email"
)

var (
	ErrCodeAlreadyExists = errors.New(codes.DuplicateError, StatusCodeAlreadyExists)
	ErrCodeNotFound      = errors.New(codes.NotFound, StatusCodeNotFound)
	ErrInvalidCode       = errors.New(codes.BadRequest, StatusInvalidCode)
	ErrInvalidEmail      = errors.New(codes.BadRequest, StatusInvalidEmail)
)

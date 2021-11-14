package user

import (
	"nnw_s/pkg/codes"

	"nnw_s/pkg/errors"
)

const (
	StatusUserNotFound      errors.Status = "user_not_found"
	StatusInvalidEmail      errors.Status = "invalid_user_email"
	StatusInvalidPassword   errors.Status = "invalid_user_password"
	StatusUserAlreadyExists errors.Status = "user_already_exists"
	StatusInvalidRequest    errors.Status = "invalid_request"
)

var (
	ErrNotFound        = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail    = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidPassword = errors.New(codes.BadRequest, StatusInvalidPassword)
	ErrAlreadyExists   = errors.New(codes.DuplicateError, StatusUserAlreadyExists)
	ErrInvalidRequest  = errors.New(codes.BadRequest, StatusInvalidRequest)
)

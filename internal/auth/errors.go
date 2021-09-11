package auth

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusUserNotFound      errors.Status = "user_not_found"
	StatusInvalidCode       errors.Status = "invalid_code"
	StatusInvalidData       errors.Status = "invalid_data"
	StatusUserAlreadyExists errors.Status = "user_already_exists"
	StatusInvalidJson       errors.Status = "invalid_json"
	StatusWrongToken        errors.Status = "wrong_token"
)

var (
	ErrNotFound      = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidCode   = errors.New(codes.BadRequest, StatusInvalidCode)
	ErrInvalidData   = errors.New(codes.BadRequest, StatusInvalidData)
	ErrAlreadyExists = errors.New(codes.DuplicateError, StatusUserAlreadyExists)
	ErrInvalidJson   = errors.New(codes.BadRequest, StatusInvalidJson)
	ErrWrongToken    = errors.New(codes.BadRequest, StatusWrongToken)
)

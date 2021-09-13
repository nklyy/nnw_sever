package jwt

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusTokenNotFound       errors.Status = "token_not_found"
	StatusTokenAlreadyExists  errors.Status = "token_already_exists"
	StatusTokenDoesNotValid   errors.Status = "token_invalid"
	StatusTokenHasBeenExpired errors.Status = "token_expired"
)

var (
	ErrTokenDoesNotValid   = errors.New(codes.Unauthorized, StatusTokenDoesNotValid)
	ErrTokenHasBeenExpired = errors.New(codes.Unauthorized, StatusTokenHasBeenExpired)
	ErrNotFound            = errors.New(codes.NotFound, StatusTokenNotFound)
	ErrAlreadyExists       = errors.New(codes.DuplicateError, StatusTokenAlreadyExists)
)

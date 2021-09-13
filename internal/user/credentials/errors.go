package credentials

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusPasswordsDoesNotMatch errors.Status = "passwords_does_not_match"
)

var (
	ErrInvalidPassword = errors.New(codes.Unauthorized, StatusPasswordsDoesNotMatch)
)

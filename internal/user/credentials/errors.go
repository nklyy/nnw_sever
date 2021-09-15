package credentials

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidPassword errors.Status = "invalid_password"
)

var (
	ErrInvalidPassword = errors.New(codes.Unauthorized, StatusInvalidPassword)
)

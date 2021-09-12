package mfa

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidMFACode errors.Status = "invalid_mfa_code"
)

var (
	ErrInvalidMFACode = errors.New(codes.Unauthorized, StatusInvalidMFACode)
)

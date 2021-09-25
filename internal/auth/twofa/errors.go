package twofa

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidTwoFACode errors.Status = "invalid_two_fa_code"
)

var (
	ErrInvalidTwoFACode = errors.New(codes.Unauthorized, StatusInvalidTwoFACode)
)

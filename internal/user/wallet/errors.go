package wallet

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidRequest errors.Status = "invalid_request"
	StatusInvalidWallet  errors.Status = "invalid_wallet"
)

var (
	ErrInvalidRequest = errors.New(codes.BadRequest, StatusInvalidRequest)
	ErrInvalidWallet  = errors.New(codes.BadRequest, StatusInvalidWallet)
)

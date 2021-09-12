package auth

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidRequest errors.Status = "invalid_request"

	StatusInvalidCode errors.Status = "invalid_code"
	StatusInvalidData errors.Status = "invalid_data"
	StatusInvalidJson errors.Status = "invalid_json"
	StatusWrongToken  errors.Status = "wrong_token"
)

var (
	ErrInvalidRequest = errors.New(codes.BadRequest, StatusInvalidRequest)

	ErrInvalidCode = errors.New(codes.BadRequest, StatusInvalidCode)
	ErrInvalidData = errors.New(codes.BadRequest, StatusInvalidData)
	ErrInvalidJson = errors.New(codes.BadRequest, StatusInvalidJson)
	ErrWrongToken  = errors.New(codes.BadRequest, StatusWrongToken)
)

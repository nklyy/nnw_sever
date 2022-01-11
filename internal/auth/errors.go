package auth

import (
	"nnw_s/pkg/codes"
	"nnw_s/pkg/errors"
)

const (
	StatusInvalidRequest           errors.Status = "invalid_request"
	StatusPermissionDenied         errors.Status = "permission_denied"
	StatusUnauthorized             errors.Status = "unauthorized"
	StatusInvalidDTO               errors.Status = "invalid_dto"
	StatusFailedCreateCode         errors.Status = "failed_create_code"
	StatusFailedSendEmail          errors.Status = "failed_send_email"
	StatusFailedGenerateTwoFaImage errors.Status = "failed_generate_twoFa_image"
	StatusInvalidCode              errors.Status = "invalid_code"
	StatusInvalidData              errors.Status = "invalid_data"
	StatusInvalidJson              errors.Status = "invalid_json"
	StatusWrongToken               errors.Status = "wrong_token"
)

var (
	ErrInvalidRequest           = errors.New(codes.BadRequest, StatusInvalidRequest)
	ErrPermissionDenied         = errors.New(codes.Forbidden, StatusPermissionDenied)
	ErrUnauthorized             = errors.New(codes.Unauthorized, StatusUnauthorized)
	ErrInvalidDTO               = errors.New(codes.InternalError, StatusInvalidDTO)
	ErrFailedCreateCode         = errors.New(codes.InternalError, StatusFailedCreateCode)
	ErrFailedSendEmail          = errors.New(codes.InternalError, StatusFailedSendEmail)
	ErrFailedGenerateTwoFaImage = errors.New(codes.InternalError, StatusFailedGenerateTwoFaImage)
	ErrInvalidCode              = errors.New(codes.InternalError, StatusInvalidCode)
)

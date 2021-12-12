package user

import (
	"nnw_s/pkg/codes"

	"nnw_s/pkg/errors"
)

const (
	StatusUserNotFound               errors.Status = "user_not_found"
	StatusInvalidEmail               errors.Status = "invalid_user_email"
	StatusInvalidPassword            errors.Status = "invalid_user_password"
	StatusUserAlreadyExists          errors.Status = "user_already_exists"
	StatusInvalidRequest             errors.Status = "invalid_request"
	StatusFailedUpdateUser           errors.Status = "failed_update_user"
	StatusUserAlreadyVerify          errors.Status = "already_verify"
	StatusUserAlreadyActive          errors.Status = "already_active"
	StatusUserAlreadyVerifyAndActive errors.Status = "already_verify_and_active"
	StatusUserDoesNotVerify          errors.Status = "user_does_not_verify"
	StatusUserDoesNotActive          errors.Status = "user_does_not_active"
)

var (
	ErrNotFound                     = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail                 = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidPassword              = errors.New(codes.BadRequest, StatusInvalidPassword)
	ErrAlreadyExists                = errors.New(codes.DuplicateError, StatusUserAlreadyExists)
	ErrInvalidRequest               = errors.New(codes.BadRequest, StatusInvalidRequest)
	ErrUserAlreadyVerify            = errors.New(codes.BadRequest, StatusUserAlreadyVerify)
	ErrUserAlreadyActive            = errors.New(codes.BadRequest, StatusUserAlreadyActive)
	ErrUserAlreadyActiveAndVerified = errors.New(codes.BadRequest, StatusUserAlreadyVerifyAndActive)
	ErrUserDoesNotVerify            = errors.New(codes.Forbidden, StatusUserDoesNotVerify)
	ErrUserDoesNotActive            = errors.New(codes.Forbidden, StatusUserDoesNotActive)
	ErrFailedUpdateUser             = errors.New(codes.InternalError, StatusFailedUpdateUser)
)

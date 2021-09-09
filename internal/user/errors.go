package user

import (
	"nnw_s/pkg/codes"

	"nnw_s/pkg/errors"
)

const (
	StatusUserNotFound         errors.Status = "user_not_found"
	StatusInvalidEmail         errors.Status = "invalid_user_email"
	StatusInvalidPassword      errors.Status = "invalid_user_password"
	StatusTemplateDataNotFound errors.Status = "template_data_not_found"
)

var (
	ErrUserNotFound         = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail         = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidPassword      = errors.New(codes.BadRequest, StatusInvalidPassword)
	ErrTemplateDataNotFound = errors.New(codes.NotFound, StatusTemplateDataNotFound)
)

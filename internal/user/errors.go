package user

import (
	"nnw_s/pkg/codes"

	"nnw_s/pkg/errors"
)

const (
	StatusUserNotFound         errors.Status = "user_not_found"
	StatusInvalidEmail         errors.Status = "invalid_user_email"
	StatusTemplateDataNotFound errors.Status = "template_data_not_found"
)

var (
	ErrUserNotFound         = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail         = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrTemplateDataNotFound = errors.New(codes.NotFound, StatusTemplateDataNotFound)
)

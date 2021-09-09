package user

import (
	"nnw_s/pkg/codes"

	"nnw_s/pkg/errors"
)

const (
	StatusUserNotFound         errors.Status = "user_not_found"
	StatusTemplateDataNotFound errors.Status = "template_data_not_found"
)

var (
	ErrUserNotFound         = errors.New(codes.NotFound, StatusUserNotFound)
	ErrTemplateDataNotFound = errors.New(codes.NotFound, StatusTemplateDataNotFound)
)

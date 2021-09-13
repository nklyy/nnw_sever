package auth

import (
	"nnw_s/pkg/errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type BaseValidator struct{}

func (v *BaseValidator) Validate() error {
	if err := validator.New().Struct(v); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.WithMessage(ErrInvalidRequest, err.Error())
		}

		validationErr := ErrInvalidRequest
		for _, err := range err.(validator.ValidationErrors) {
			validationErr = errors.WithMessage(validationErr, err.Error())
		}
		return validationErr
	}
	return nil
}

type RegisterUserDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`

	BaseValidator
}

type VerifyUserDTO struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`

	BaseValidator
}

type SetupMfaDTO struct {
	Email string `json:"email" validate:"required"`
	BaseValidator
}

type ActivateUserDTO struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`

	BaseValidator
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`

	BaseValidator
}

type TokenDTO struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expired_at"`
}

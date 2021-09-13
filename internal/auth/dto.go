package auth

import (
	"nnw_s/pkg/errors"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const passwordMinLength = 10

type BaseValidator struct{}

func (v *BaseValidator) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", v.validatePassword)

	if err := validate.Struct(v); err != nil {
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

func (v *BaseValidator) validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < passwordMinLength {
		return false
	}

	var (
		containsUpper bool
		containsLower bool
		containsDigit bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			containsUpper = true
		} else if unicode.IsLower(char) {
			containsLower = true
		} else if unicode.IsDigit(char) {
			containsDigit = true
		}
	}
	return containsUpper && containsLower && containsDigit
}

type RegisterUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	BaseValidator
}

type VerifyUserDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`

	BaseValidator
}

type ResendActivationEmailDTO struct {
	Email string `json:"email" validate:"required,email"`

	BaseValidator
}

type SetupMfaDTO struct {
	Email string `json:"email" validate:"required,email"`

	BaseValidator
}

type ActivateUserDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`

	BaseValidator
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required,len=6,numeric"`

	BaseValidator
}

type TokenDTO struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expired_at"`
}

package auth

import (
	"nnw_s/pkg/errors"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

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

	// minimum eight and maximum 10 characters, at least one uppercase letter, one lowercase letter, one number and one special character
	regex, _ := regexp.Compile("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,10}$")
	result := regex.MatchString(password)
	return !result
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

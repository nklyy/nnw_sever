package auth

import (
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const passwordMinLength = 8

func Validate(dto interface{}, shift int) error {
	validate := validator.New()
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		decodedPassword, err := helpers.CaesarShift(password, -shift)
		if err != nil {
			return false
		}

		if len(decodedPassword) >= passwordMinLength {
			return true
		}

		var (
			containsUpper bool
			containsLower bool
			containsDigit bool
		)

		for _, char := range decodedPassword {
			if unicode.IsUpper(char) {
				containsUpper = true
			} else if unicode.IsLower(char) {
				containsLower = true
			} else if unicode.IsDigit(char) {
				containsDigit = true
			}
		}
		return containsUpper && containsLower && containsDigit
	})

	if err := validate.Struct(dto); err != nil {
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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type VerifyUserDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type ResendActivationEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type SetupTwoFaDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ActivateUserDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginCodeDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LogoutCodeDTO struct {
	Token string `json:"token" validate:"required"`
}

type TokenDTO struct {
	Token    string    `json:"token" validate:"required"`
	ExpireAt time.Time `json:"expired_at" validate:"required"`
}

type ValidateTokenDTO struct {
	Token string `json:"token" validate:"required"`
}

type ResetPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResendResetPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordCodedDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type SetupNewPasswordDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

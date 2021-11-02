package wallet

import (
	"github.com/go-playground/validator/v10"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"
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

		if len(decodedPassword) < passwordMinLength {
			return false
		}

		return true
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

type CreateWalletDTO struct {
	Password string `json:"password" validate:"required,password"`
	Backup   *bool  `json:"backup" validate:"required"`
	Jwt      string `json:"jwt" validate:"required"`
}

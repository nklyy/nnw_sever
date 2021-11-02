package user

import (
	"github.com/go-playground/validator/v10"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/wallet"
	"time"
)

func Validate(dto interface{}) error {
	validate := validator.New()
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

type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserDTO struct {
	Jwt string `json:"jwt" validate:"required"`
}

type DTO struct {
	ID         string            `json:"id"`
	Email      string            `json:"email"`
	Password   string            `json:"password"`
	SecretOTP  string            `json:"secret_otp"`
	Status     string            `json:"status"`
	Wallet     *[]*wallet.Wallet `json:"wallet"`
	IsVerified bool              `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserResponseDTO struct {
	Email      string            `json:"email"`
	Status     string            `json:"status"`
	IsVerified bool              `json:"is_verified"`
	Wallet     *[]*wallet.Wallet `json:"wallet"`
}

func NormalizeGetUserResponseDTO(dto *DTO) *GetUserResponseDTO {
	return &GetUserResponseDTO{
		Email:      dto.Email,
		Status:     dto.Status,
		IsVerified: dto.IsVerified,
		Wallet:     dto.Wallet,
	}
}

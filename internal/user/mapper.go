package user

import (
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapToDTO(u *User) *DTO {
	var secretOTP string
	if u.Credentials.SecretOTP != nil {
		secretOTP = *u.Credentials.SecretOTP
	}
	return &DTO{
		ID:         u.ID.Hex(),
		Email:      u.Email,
		Password:   u.Credentials.Password,
		SecretOTP:  secretOTP,
		Status:     string(u.Status),
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func MapToEntity(dto *DTO) (*User, error) {
	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	return &User{
		ID:    id,
		Email: dto.Email,
		Credentials: &credentials.Credentials{
			Password:  dto.Password,
			SecretOTP: &dto.SecretOTP,
		},
		Status:     Status(dto.Status),
		IsVerified: dto.IsVerified,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}, nil
}

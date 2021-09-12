package user

import "go.mongodb.org/mongo-driver/bson/primitive"

func MapToDTO(u *User) *DTO {
	return &DTO{
		ID:         u.ID.Hex(),
		Email:      u.Email,
		Password:   u.Credentials.Password,
		SecretOTP:  *u.Credentials.SecretOTP,
		Status:     string(u.Status),
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func MapToEntity(dto *DTO) *User {
	// todo: can't convert string to primitive.ObjectID
	return &User{
		ID: primitive.ObjectID([12]byte(dto.ID)),
	}
}

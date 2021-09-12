package user

import (
	"nnw_s/pkg/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	Status       string             `bson:"status"`
	VerifyEmail  bool               `bson:"verify_email"`
	SecretOTPKey *string            `bson:"secret_otp_key"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

func NewUser(email string, password string, secretOTP *string) (*User, error) {
	// put other validation
	if email == "" {
		return nil, errors.WithMessage(ErrInvalidEmail, "should be not empty")
	}
	if password == "" {
		return nil, errors.WithMessage(ErrInvalidPassword, "should be not empty")
	}

	return &User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		Password:     password,
		Status:       "disable",
		VerifyEmail:  false,
		SecretOTPKey: secretOTP,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func NewDisableUser(email string, password string, oldUser *User) (*User, error) {
	// put other validation
	if email == "" {
		return nil, errors.WithMessage(ErrInvalidEmail, "should be not empty")
	}
	if password == "" {
		return nil, errors.WithMessage(ErrInvalidPassword, "should be not empty")
	}

	return &User{
		ID:           oldUser.ID,
		Email:        email,
		Password:     password,
		Status:       oldUser.Status,
		VerifyEmail:  oldUser.VerifyEmail,
		SecretOTPKey: oldUser.SecretOTPKey,
		CreatedAt:    oldUser.CreatedAt,
		UpdatedAt:    time.Now(),
	}, nil
}

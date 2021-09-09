package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	SecretOTPKey string             `bson:"secret_otp_key"`
	VerifyEmail  bool               `bson:"verify_email"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

func NewUser(email, passwordHash, secretOTP string) (*User, error) {
	// put other validation
	if email == "" {
		return nil, ErrInvalidEmail
	}

	return &User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		Password:     passwordHash,
		SecretOTPKey: secretOTP,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

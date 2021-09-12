package verification

import (
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Code struct {
	ID        primitive.ObjectID `bson:"_id"`
	Code      string             `bson:"code"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewCode(email string) (*Code, error) {
	if email == "" {
		return nil, errors.WithMessage(ErrInvalidEmail, "email cannot be empty")
	}
	return &Code{
		ID:        primitive.NewObjectID(),
		Code:      helpers.EmailCode(),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

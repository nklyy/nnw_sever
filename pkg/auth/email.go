package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Email struct {
	ID        primitive.ObjectID `bson:"_id"`
	Code      string             `bson:"code"`
	Email     string             `bson:"email"`
	EmailType string             `bson:"email_type"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

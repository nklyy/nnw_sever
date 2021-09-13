package jwt

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWT struct {
	ID        primitive.ObjectID `bson:"_id"`
	Jwt       string             `bson:"jwt"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewJWT(token string) *JWT {
	return &JWT{
		ID:        primitive.NewObjectID(),
		Jwt:       token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

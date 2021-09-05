package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JWTData struct {
	ID        primitive.ObjectID `bson:"_id"`
	Jwt       string             `bson:"jwt"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

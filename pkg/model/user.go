package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Login        string             `bson:"login"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	SecretOTPKey string             `bson:"secret_otp_key"`
}

type TemplateData struct {
	ID        primitive.ObjectID `bson:"_id"`
	Uid       string             `bson:"uid"`
	TwoFAS    string             `bson:"two_fas"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

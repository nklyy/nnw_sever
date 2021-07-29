package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `bson:"id"`
	Login        string             `bson:"login"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	SecretOTPKey string             `bson:"secret_otp_key"`
}

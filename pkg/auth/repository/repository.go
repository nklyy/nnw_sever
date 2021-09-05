package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"nnw_s/config"
	model2 "nnw_s/pkg/auth/model"
)

type Authorization interface {
	GetJwtDb(id string) (*string, error)
	CreateJwtDb(jwtData model2.JWTData) (string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database, cfg config.Configurations) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db, cfg),
	}
}

package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"nnw_s/pkg/model"
)

type Authorization interface {
	GetUserByIdDb(userId string) (*model.User, error)
	GetUserByLoginDb(login string) (*model.User, error)
	GetTemplateUserDataByIdDb(uid string) (*model.TemplateData, error)
	CreateUserDb(user model.User) (*string, error)
	CreateTemplateUserDataDb(templateData model.TemplateData) (*string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db),
	}
}

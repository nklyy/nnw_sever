package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"nnw_s/config"
	model2 "nnw_s/pkg/user/model"
)

type User interface {
	GetUserByIdDb(userId string) (*model2.User, error)
	GetUserByEmailDb(email string) (*model2.User, error)
	GetTemplateUserDataByIdDb(uid string) (*model2.TemplateData, error)

	CreateUserDb(user model2.User) (*string, error)
	CreateTemplateUserDataDb(templateData model2.TemplateData) (*string, error)
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Database, cfg config.Configurations) *Repository {
	return &Repository{
		User: NewUserMongo(db, cfg),
	}
}

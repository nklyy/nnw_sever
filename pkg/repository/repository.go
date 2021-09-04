package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"nnw_s/config"
	"nnw_s/pkg/model"
)

type Authorization interface {
	GetUserByIdDb(userId string) (*model.User, error)
	GetUserByEmailDb(email string) (*model.User, error)
	GetTemplateUserDataByIdDb(uid string) (*model.TemplateData, error)
	GetJwtDb(id string) (*string, error)
	GetEmailDb(email string, code string, emailType string) (*model.Email, error)

	CreateUserDb(user model.User) (*string, error)
	CreateTemplateUserDataDb(templateData model.TemplateData) (*string, error)
	CreateJwtDb(jwtData model.JWTData) (string, error)
	CreateEmailDb(emailData model.Email) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database, cfg config.Configurations) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db, cfg),
	}
}

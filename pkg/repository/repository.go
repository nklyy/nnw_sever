package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"nnw_s/pkg/model"
)

type User interface {
	GetUserByIdDb(userId string) (*model.User, error)
	CreateUserDb(user model.User) (*string, error)
	UpdateUserDb(user model.User) error
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User: NewUserMongo(db),
	}
}
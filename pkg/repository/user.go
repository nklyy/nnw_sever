package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"nnw_s/pkg/model"
)

type UserMongo struct {
	db *mongo.Database
}

func NewUserMongo(db *mongo.Database) *UserMongo {
	return &UserMongo{
		db: db,
	}
}

func (ur *UserMongo) GetUserByIdDb(userId string) (*model.User, error) {
	return nil, nil
}

func (ur *UserMongo) CreateUserDb(user model.User) (*string, error) {
	return nil, nil
}

func (ur *UserMongo) UpdateUserDb(user model.User) error {
	return nil
}

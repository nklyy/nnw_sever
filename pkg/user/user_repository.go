package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nnw_s/config"
)

type IUserRepository interface {
	GetUserByIdDb(userId string) (*User, error)
	GetUserByEmailDb(email string) (*User, error)
	GetTemplateUserDataByIdDb(uid string) (*TemplateData, error)

	CreateUserDb(user User) (*string, error)
	CreateTemplateUserDataDb(templateData TemplateData) (*string, error)
}

type UserRepository struct {
	db  *mongo.Database
	cfg config.Configurations
}

func NewUserRepository(db *mongo.Database, cfg config.Configurations) *UserRepository {
	return &UserRepository{
		db:  db,
		cfg: cfg,
	}
}

func (ar *UserRepository) GetUserByIdDb(userId string) (*User, error) {
	return nil, nil
}

func (ar *UserRepository) GetUserByEmailDb(email string) (*User, error) {
	var user User

	err := ar.db.Collection("user").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *UserRepository) GetTemplateUserDataByIdDb(uid string) (*TemplateData, error) {
	var templateUser TemplateData

	err := ar.db.Collection("templateUserData").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&templateUser)
	if err != nil {
		return nil, err
	}

	return &templateUser, nil
}

func (ar *UserRepository) CreateUserDb(user User) (*string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := ar.db.Collection("user").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return nil, err
	}

	_, err = ar.db.Collection("user").InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ar *UserRepository) CreateTemplateUserDataDb(templateData TemplateData) (*string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(500),
	}

	_, err := ar.db.Collection("templateUserData").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return nil, err
	}

	_, err = ar.db.Collection("templateUserData").InsertOne(context.TODO(), templateData)
	if err != nil {
		return nil, err
	}

	return &templateData.Uid, nil
}

package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nnw_s/config"
	"nnw_s/pkg/user/model"
)

type UserMongo struct {
	db  *mongo.Database
	cfg config.Configurations
}

func NewUserMongo(db *mongo.Database, cfg config.Configurations) *UserMongo {
	return &UserMongo{
		db:  db,
		cfg: cfg,
	}
}

func (ar *UserMongo) GetUserByIdDb(userId string) (*model.User, error) {
	return nil, nil
}

func (ar *UserMongo) GetUserByEmailDb(email string) (*model.User, error) {
	var user model.User

	err := ar.db.Collection("user").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *UserMongo) GetTemplateUserDataByIdDb(uid string) (*model.TemplateData, error) {
	var templateUser model.TemplateData

	err := ar.db.Collection("templateUserData").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&templateUser)
	if err != nil {
		return nil, err
	}

	return &templateUser, nil
}

func (ar *UserMongo) CreateUserDb(user model.User) (*string, error) {
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

func (ar *UserMongo) CreateTemplateUserDataDb(templateData model.TemplateData) (*string, error) {
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

package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nnw_s/config"
	"nnw_s/pkg/model"
)

type AuthMongo struct {
	db  *mongo.Database
	cfg config.Configurations
}

func NewAuthMongo(db *mongo.Database, cfg config.Configurations) *AuthMongo {
	return &AuthMongo{
		db:  db,
		cfg: cfg,
	}
}

func (ar *AuthMongo) GetUserByIdDb(userId string) (*model.User, error) {
	return nil, nil
}

func (ar *AuthMongo) GetUserByLoginDb(login string) (*model.User, error) {
	var user model.User

	err := ar.db.Collection("user").FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *AuthMongo) GetTemplateUserDataByIdDb(uid string) (*model.TemplateData, error) {
	var templateUser model.TemplateData

	err := ar.db.Collection("templateUserData").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&templateUser)
	if err != nil {
		return nil, err
	}

	return &templateUser, nil
}

func (ar *AuthMongo) CreateUserDb(user model.User) (*string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"login": 1}, // index in ascending order or -1 for descending order
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

func (ar *AuthMongo) CreateTemplateUserDataDb(templateData model.TemplateData) (*string, error) {
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
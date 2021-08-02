package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (ur *UserMongo) GetUserByLoginDb(login string) (*model.User, error) {
	var user model.User

	err := ur.db.Collection("user").FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserMongo) GetTemplateUserDataByIdDb(uid string) (*model.TemplateData, error) {
	var templateUser model.TemplateData

	err := ur.db.Collection("templateUserData").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&templateUser)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &templateUser, nil
}

func (ur *UserMongo) CreateUserDb(user model.User) (*string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"login": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := ur.db.Collection("user").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return nil, err
	}

	_, err = ur.db.Collection("user").InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur *UserMongo) CreateTemplateUserDataDb(templateData model.TemplateData) (*string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(500),
	}

	_, err := ur.db.Collection("templateUserData").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return nil, err
	}

	_, err = ur.db.Collection("templateUserData").InsertOne(context.TODO(), templateData)
	if err != nil {
		return nil, err
	}

	return &templateData.Uid, nil
}

func (ur *UserMongo) UpdateUserDb(user model.User) error {
	return nil
}

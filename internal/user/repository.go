package user

import (
	"context"
	"nnw_s/config"
	"nnw_s/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const templateDataExpiry = 500

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetTemplateUserDataByID(ctx context.Context, uid string) (*TemplateData, error)

	CreateUser(ctx context.Context, user User) (string, error)
	CreateTemplateUserData(ctx context.Context, templateData TemplateData) (string, error)
}

type repository struct {
	db  *mongo.Database
	cfg config.Config
}

func NewRepository(db *mongo.Database, cfg config.Config) *repository {
	return &repository{
		db:  db,
		cfg: cfg,
	}
}

func (repo *repository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return nil, nil
}

func (repo *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := repo.db.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, errors.NewInternal(err.Error())
	}

	return &user, nil
}

func (repo *repository) GetTemplateUserDataByID(ctx context.Context, uid string) (*TemplateData, error) {
	var templateUser TemplateData

	err := repo.db.Collection("templateUserData").FindOne(ctx, bson.M{"uid": uid}).Decode(&templateUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTemplateDataNotFound
		}
		return nil, errors.NewInternal(err.Error())
	}

	return &templateUser, nil
}

func (repo *repository) CreateUser(ctx context.Context, user User) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.db.Collection("user").Indexes().CreateOne(ctx, mod)
	if err != nil {
		return "", errors.NewInternal(err.Error())
	}

	_, err = repo.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		return "", errors.NewInternal(err.Error())
	}
	return user.ID.String(), nil
}

func (repo *repository) CreateTemplateUserData(ctx context.Context, templateData TemplateData) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(templateDataExpiry),
	}

	_, err := repo.db.Collection("templateUserData").Indexes().CreateOne(ctx, mod)
	if err != nil {
		return "", errors.NewInternal(err.Error())
	}

	_, err = repo.db.Collection("templateUserData").InsertOne(ctx, templateData)
	if err != nil {
		return "", errors.NewInternal(err.Error())
	}

	return templateData.Uid, nil
}

package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nnw_s/config"
)

const (
	emailExpiry = 600
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetJwt(ctx context.Context, id string) (string, error)
	GetEmail(ctx context.Context, email, code, emailType string) (*Email, error)
	CreateJwt(ctx context.Context, jwtData JWT) (string, error)
	CreateEmail(ctx context.Context, emailData Email) error
}

type repository struct {
	db  *mongo.Database
	cfg config.Config
}

func NewRepository(db *mongo.Database, cfg config.Config) Repository {
	return &repository{
		db:  db,
		cfg: cfg,
	}
}

func (repo *repository) GetJwt(ctx context.Context, id string) (string, error) {
	var jwtData JWT

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	err = repo.db.Collection("jwt").FindOne(ctx, bson.M{"_id": objectId}).Decode(&jwtData)
	if err != nil {
		return "", err
	}

	return jwtData.Jwt, nil
}

func (repo *repository) GetEmail(ctx context.Context, email, code, emailType string) (*Email, error) {
	var emailData Email

	err := repo.db.Collection("email").FindOne(ctx, bson.M{"email": email, "code": code, "email_type": emailType}).Decode(&emailData)

	if err != nil {
		return nil, err
	}

	return &emailData, nil
}

func (repo *repository) CreateJwt(ctx context.Context, jwtData JWT) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(60),
	}

	_, err := repo.db.Collection("jwt").Indexes().CreateOne(ctx, mod)
	if err != nil {
		return "", err
	}

	_, err = repo.db.Collection("jwt").InsertOne(ctx, jwtData)
	if err != nil {
		return "", err
	}

	return jwtData.ID.Hex(), nil
}

func (repo *repository) CreateEmail(ctx context.Context, emailData Email) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(emailExpiry),
	}

	_, err := repo.db.Collection("email").Indexes().CreateOne(ctx, mod)
	if err != nil {
		return err
	}

	_, err = repo.db.Collection("email").InsertOne(ctx, emailData)
	if err != nil {
		return err
	}

	return nil
}

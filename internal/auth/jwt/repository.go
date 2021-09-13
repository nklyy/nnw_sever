package jwt

import (
	"context"
	"nnw_s/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	GetJWT(ctx context.Context, id string) (*JWT, error)
	SaveJWT(ctx context.Context, jwt *JWT) (string, error)
}

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) (Repository, error) {
	if db == nil {
		return nil, errors.NewInternal("invalid db")
	}
	return &repository{db: db}, nil
}

func (repo *repository) GetJWT(ctx context.Context, id string) (*JWT, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	var jwtData JWT
	if err = repo.db.Collection("jwt").FindOne(ctx, bson.M{"_id": objectId}).Decode(&jwtData); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, errors.NewInternal(err.Error())
	}
	return &jwtData, nil
}

func (repo *repository) SaveJWT(ctx context.Context, jwt *JWT) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(60),
	}

	_, err := repo.db.Collection("jwt").Indexes().CreateOne(ctx, mod)
	if err != nil {
		return "", errors.NewInternal(err.Error())
	}

	_, err = repo.db.Collection("jwt").InsertOne(ctx, jwt)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrAlreadyExists
		}
		return "", errors.NewInternal(err.Error())
	}
	return jwt.ID.Hex(), nil
}

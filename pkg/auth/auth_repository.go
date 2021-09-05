package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nnw_s/config"
)

type IAuthRepository interface {
	GetJwtDb(id string) (*string, error)
	CreateJwtDb(jwtData JWTData) (string, error)
}

type AuthRepository struct {
	db  *mongo.Database
	cfg config.Configurations
}

func NewAuthRepository(db *mongo.Database, cfg config.Configurations) *AuthRepository {
	return &AuthRepository{
		db:  db,
		cfg: cfg,
	}
}

func (ar *AuthRepository) GetJwtDb(id string) (*string, error) {
	var jwtData JWTData

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = ar.db.Collection("jwt").FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&jwtData)
	if err != nil {
		return nil, err
	}

	return &jwtData.Jwt, nil
}

func (ar *AuthRepository) CreateJwtDb(jwtData JWTData) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(60),
	}

	_, err := ar.db.Collection("jwt").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return "", err
	}

	_, err = ar.db.Collection("jwt").InsertOne(context.TODO(), jwtData)
	if err != nil {
		return "", err
	}

	return jwtData.ID.Hex(), nil
}

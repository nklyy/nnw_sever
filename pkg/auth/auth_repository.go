package auth

import (
	"context"
	"github.com/golang/mock/mockgen/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nnw_s/config"
)

type IAuthRepository interface {
	GetJwtDb(id string) (*string, error)
	GetEmailDb(email string, code string, emailType string) (*Email, error)
	CreateJwtDb(jwtData JWTData) (string, error)
	CreateEmailDb(emailData Email) error
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

func (ar *AuthRepository) GetEmailDb(email string, code string, emailType string) (*Email, error) {
	var emailData Email

	err := ar.db.Collection("email").FindOne(context.TODO(), bson.M{"email": email, "code": code, "email_type": emailType}).Decode(&emailData)

	if err != nil {
		return nil, err
	}

	return &emailData, nil
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

func (ar *AuthRepository) CreateTemplateUserDataDb(templateData model.TemplateData) (*string, error) {
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

func (ar *AuthRepository) CreateEmailDb(emailData Email) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(600),
	}

	_, err := ar.db.Collection("email").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return err
	}

	_, err = ar.db.Collection("email").InsertOne(context.TODO(), emailData)
	if err != nil {
		return err
	}

	return nil
}

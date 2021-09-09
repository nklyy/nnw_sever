package user

import (
	"context"
	"nnw_s/pkg/errors"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	SaveUser(ctx context.Context, user *User) (string, error)
}

type repository struct {
	db  *mongo.Database
	log *logrus.Logger
}

func NewRepository(db *mongo.Database, log *logrus.Logger) *repository {
	return &repository{
		db:  db,
		log: log,
	}
}

func (repo *repository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := repo.db.Collection("user").FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			repo.log.WithContext(ctx).Errorf("unable to find user by id '%s': %v", userID, err)
			return nil, ErrUserNotFound
		}

		repo.log.WithContext(ctx).Errorf("unable to find user due to internal error: %v; id: %s", err, userID)
		return nil, errors.NewInternal(err.Error())
	}

	return &user, nil
}

func (repo *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := repo.db.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			repo.log.WithContext(ctx).Errorf("unable to find user by email '%s': %v", email, err)
			return nil, ErrUserNotFound
		}

		repo.log.WithContext(ctx).Errorf("unable to find user due to internal error: %v; email: %s", err, email)
		return nil, errors.NewInternal(err.Error())
	}

	return &user, nil
}

func (repo *repository) SaveUser(ctx context.Context, user User) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.db.Collection("user").Indexes().CreateOne(ctx, mod)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to create user index: %v", err)
		return "", errors.NewInternal(err.Error())
	}

	_, err = repo.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to insert user data to db: %v", err)
		return "", errors.NewInternal(err.Error())
	}
	return user.ID.String(), nil
}

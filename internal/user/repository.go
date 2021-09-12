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
	GetUserByEmail(ctx context.Context, email string, status string) (*User, error)

	SaveUser(ctx context.Context, user *User) (string, error)

	UpdateUser(ctx context.Context, user *User) error
	UpdateDisableUser(ctx context.Context, user *User) error
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
			return nil, ErrNotFound
		}

		repo.log.WithContext(ctx).Errorf("unable to find user due to internal error: %v; id: %s", err, userID)
		return nil, errors.NewInternal(err.Error())
	}

	return &user, nil
}

func (repo *repository) GetUserByEmail(ctx context.Context, email string, status string) (*User, error) {
	var user User
	if err := repo.db.Collection("user").FindOne(ctx, bson.M{"email": email, "status": status}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			repo.log.WithContext(ctx).Errorf("unable to find user by email '%s': %v", email, err)
			return nil, ErrNotFound
		}

		repo.log.WithContext(ctx).Errorf("unable to find user due to internal error: %v; email: %s", err, email)
		return nil, errors.NewInternal(err.Error())
	}

	return &user, nil
}

func (repo *repository) SaveUser(ctx context.Context, user *User) (string, error) {
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
		if mongo.IsDuplicateKeyError(err) {
			repo.log.WithContext(ctx).Errorf("failed to insert user data to db due to duplicate error: %v", err)
			return "", ErrAlreadyExists
		}

		repo.log.WithContext(ctx).Errorf("failed to insert user data to db: %v", err)
		return "", errors.NewInternal(err.Error())
	}
	return user.ID.Hex(), nil
}

func (repo *repository) UpdateUser(ctx context.Context, user *User) error {
	_, err := repo.db.Collection("user").UpdateOne(ctx, bson.M{"email": user.Email}, bson.D{{"$set", user}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) UpdateDisableUser(ctx context.Context, user *User) error {
	_, err := repo.db.Collection("user").UpdateOne(ctx, bson.M{"email": user.Email}, bson.D{{"$set", user}})
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to update user data in db: %v", err)
		return errors.NewInternal(err.Error())
	}

	return nil
}

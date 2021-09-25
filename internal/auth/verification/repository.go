package verification

import (
	"context"
	"nnw_s/pkg/errors"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	verificationCodeExpiry  = 600
	resetPasswordCodeExpiry = 300
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetVerificationCode(ctx context.Context, email, code string) (*Code, error)
	SaveVerificationCode(ctx context.Context, code *Code) error

	GetResetPasswordCode(ctx context.Context, email, code string) (*Code, error)
	SaveResetPasswordCode(ctx context.Context, code *Code) error
}

type repository struct {
	db  *mongo.Database
	log *logrus.Logger
}

func NewRepository(db *mongo.Database, log *logrus.Logger) (Repository, error) {
	if db == nil {
		return nil, errors.NewInternal("db cannot be nil")
	}
	if log == nil {
		return nil, errors.NewInternal("logger cannot be nil")
	}
	return &repository{db: db, log: log}, nil
}

func (repo *repository) SaveVerificationCode(ctx context.Context, code *Code) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(verificationCodeExpiry),
	}

	_, err := repo.db.Collection("verification_code").Indexes().CreateOne(ctx, mod)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to create verification code index: %v", err)
		return errors.NewInternal(err.Error())
	}

	_, err = repo.db.Collection("verification_code").InsertOne(ctx, code)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to insert verification code to db: %v", err)
		return errors.NewInternal(err.Error())
	}

	return nil
}

func (repo *repository) GetVerificationCode(ctx context.Context, email, code string) (*Code, error) {
	var verificationCode Code
	err := repo.db.
		Collection("verification_code").
		FindOne(ctx, bson.M{"email": email, "code": code}).
		Decode(&verificationCode)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			repo.log.WithContext(ctx).Errorf("unable to find verification code by email and code: %v", err)
			return nil, ErrCodeNotFound
		}

		repo.log.WithContext(ctx).Errorf("unable to find verification code due to internal error: %v", err)
		return nil, errors.NewInternal(err.Error())
	}

	return &verificationCode, nil
}

func (repo *repository) SaveResetPasswordCode(ctx context.Context, code *Code) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(resetPasswordCodeExpiry),
	}

	_, err := repo.db.Collection("reset_password_code").Indexes().CreateOne(ctx, mod)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to create reset password code index: %v", err)
		return errors.NewInternal(err.Error())
	}

	_, err = repo.db.Collection("reset_password_code").InsertOne(ctx, code)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("failed to insert reset password code to db: %v", err)
		return errors.NewInternal(err.Error())
	}

	return nil
}

func (repo *repository) GetResetPasswordCode(ctx context.Context, email, code string) (*Code, error) {
	var resetPasswordCode Code
	err := repo.db.
		Collection("reset_password_code").
		FindOne(ctx, bson.M{"email": email, "code": code}).
		Decode(&resetPasswordCode)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			repo.log.WithContext(ctx).Errorf("unable to find reset password code by email and code: %v", err)
			return nil, ErrCodeNotFound
		}

		repo.log.WithContext(ctx).Errorf("unable to find reset password code due to internal error: %v", err)
		return nil, errors.NewInternal(err.Error())
	}

	return &resetPasswordCode, nil
}

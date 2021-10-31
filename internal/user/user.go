package user

import (
	"nnw_s/internal/user/credentials"
	"nnw_s/pkg/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID       `bson:"_id"`
	Email       string                   `bson:"email"`
	Credentials *credentials.Credentials `bson:"credentials"`
	Status      Status                   `bson:"status"`
	IsVerified  bool                     `bson:"is_verified"`
	BtcWallet   string                   `bson:"btc_wallet"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewUser(email string, credentials *credentials.Credentials) (*User, error) {
	if email == "" {
		return nil, errors.WithMessage(ErrInvalidEmail, "should be not empty")
	}
	return &User{
		ID:          primitive.NewObjectID(),
		Email:       email,
		Credentials: credentials,
		Status:      Disabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (u *User) IsActive() bool {
	return u.Status == Active
}

func (u *User) SetToVerified() {
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}

func (u *User) SetToActive() {
	u.Status = Active
	u.UpdatedAt = time.Now()
}

func (u *User) SetBtcWallet(walletName string) {
	u.BtcWallet = walletName
}

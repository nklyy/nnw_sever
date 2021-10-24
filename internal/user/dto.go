package user

import (
	"time"
)

type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DTO struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	SecretOTP  string `json:"secret_otp"`
	Status     string `json:"status"`
	BtcWallet  string `json:"btc_wallet"`
	IsVerified bool   `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

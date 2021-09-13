package jwt

import "time"

type DTO struct {
	ID       string    `json:"id"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

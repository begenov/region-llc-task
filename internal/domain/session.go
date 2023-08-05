package domain

import "time"

type Session struct {
	RefreshToken string    `json:"refresh_token"`
	ExpirationAt time.Time `json:"expiration_at"`
}

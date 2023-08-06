package domain

import "time"

type Session struct {
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	ExpirationAt time.Time `bson:"expiration_at" json:"expiration_at"`
}

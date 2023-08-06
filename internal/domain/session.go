package domain

import "time"

type Session struct {
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	ExpirationAt time.Time `bson:"expiration_at" json:"expiration_at"`
}

type Token struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

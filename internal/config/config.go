package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APIEndpoint string        `envconfig:"API_ENDPOINT" required:"true"`
	Mongo       ConfigMongo   `envconfig:"MONGO" required:"true"`
	Redis       ConfigRedis   `envconfig:"REDIS" required:"true"`
	Session     ConfigSession `envconfig:"SESSION" required:"true"`
}

type ConfigMongo struct {
	Uri      string `envconfig:"URI" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Name     string `envconfig:"NAME" required:"true"`
}

type ConfigSession struct {
	SignKey         string        `envconfig:"SIGN_KEY" required:"true"`
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL" required:"true"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" required:"true"`
}

type ConfigRedis struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     int    `envconfig:"PORT" required:"true"`
	Password string `envconfig:"PASSWORD" required:"false"`
	DB       int    `envconfig:"DB" required:"true"`
}

func NewConfig(fpath string) (*Config, error) {
	err := godotenv.Load(fpath)
	if err != nil {
		return nil, fmt.Errorf("godotenv.Load(): %s", err.Error())
	}

	var config Config

	err = envconfig.Process("", &config)
	if err != nil {
		return nil, fmt.Errorf("envconfig.Process(): %s", err.Error())
	}

	return &config, nil
}

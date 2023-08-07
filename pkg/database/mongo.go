package database

import (
	"context"
	"fmt"

	"github.com/begenov/region-llc-task/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, cfg config.ConfigMongo) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(cfg.Uri)
	if cfg.User != "" && cfg.Password != "" {
		opts.SetAuth(options.Credential{
			Username: cfg.User, Password: cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect(): %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("client.Ping(): %v", err)
	}

	return client, nil
}

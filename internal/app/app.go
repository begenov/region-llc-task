package app

import (
	"context"
	"fmt"
	"time"

	"github.com/begenov/region-llc-task/internal/config"
	"github.com/begenov/region-llc-task/internal/delivery/http"
	mongorepo "github.com/begenov/region-llc-task/internal/repository/mongo"
	redisrepo "github.com/begenov/region-llc-task/internal/repository/redis"
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/begenov/region-llc-task/pkg/auth"
	"github.com/begenov/region-llc-task/pkg/database"
	"github.com/begenov/region-llc-task/pkg/hash"
	"github.com/begenov/region-llc-task/pkg/redis"
)

const timeout = 10 * time.Second

func Run(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	redisClient, err := redis.CreateClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("redis.CreateClient(): %v", err)
	}
	defer redisClient.Close()

	mongoClient, err := database.NewClient(ctx, cfg.Mongo)
	if err != nil {
		return fmt.Errorf("database.NewClient(): %v", err)
	}
	db := mongoClient.Database(cfg.Mongo.Name)
	defer db.Client().Disconnect(ctx)

	hash := hash.NewHash()
	manager, err := auth.NewManager(cfg.Session.SignKey)
	if err != nil {
		return fmt.Errorf("auth.NewManager(): %v", err)
	}

	redisRepo := redisrepo.NewRedis(redisClient)
	userRepo := mongorepo.NewUserRepo(db)
	todoRepo := mongorepo.NewTodoRepo(db)

	userService := service.NewUserService(userRepo, hash, manager, cfg.Session.AccessTokenTTL,
		cfg.Session.RefreshTokenTTL, redisRepo)
	todoService := service.NewTodoService(todoRepo, userRepo)

	server := http.NewServer(userService, todoService, manager)

	if err := server.Init(cfg.APIEndpoint); err != nil {
		return fmt.Errorf("server.Init(): %v", err)
	}

	return nil
}

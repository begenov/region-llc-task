package app

import (
	"context"
	"fmt"
	"time"

	"github.com/begenov/region-llc-task/internal/config"
	"github.com/begenov/region-llc-task/internal/delivery/http"
	mongorepo "github.com/begenov/region-llc-task/internal/repository/mongo"
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/begenov/region-llc-task/pkg/database"
)

const timeout = 10 * time.Second

func Run(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	mongoClient, err := database.NewClient(ctx, cfg.Mongo)
	if err != nil {
		return fmt.Errorf("database.NewClient(): %v", err)
	}
	db := mongoClient.Database(cfg.Mongo.Name)
	defer db.Client().Disconnect(ctx)

	userRepo := mongorepo.NewUserRepo(db)
	todoRepo := mongorepo.NewTodoRepo(db)

	userService := service.NewUserService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	server := http.NewServer(userService, todoService)

	if err := server.Init(cfg.APIEndpoint); err != nil {
		return fmt.Errorf("server.Init(): %v", err)
	}

	return nil
}

package main

import (
	"github.com/begenov/region-llc-task/internal/app"
	"github.com/begenov/region-llc-task/internal/config"
	"github.com/begenov/region-llc-task/pkg/logger"
)

// @title Todo List API
// @version 1.0
// @description API Server for Todo List

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		logger.Fatalf("config.NewConfig(): %v", err)
	}

	if err := app.Run(cfg); err != nil {
		logger.Errorf("app.Run(): %v", err)
	}
}

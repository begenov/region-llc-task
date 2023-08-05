package main

import (
	"github.com/begenov/region-llc-task/internal/app"
	"github.com/begenov/region-llc-task/internal/config"
	"github.com/begenov/region-llc-task/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		logger.Fatalf("config.NewConfig(): %v", err)
	}

	if err := app.Run(cfg); err != nil {
		logger.Errorf("app.Run(): %v", err)
	}
}

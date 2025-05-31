package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/routes"
	"github.com/artnikel/vacancystats/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logging.NewLogger(cfg.Logging.Path)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	pudge, err := storage.NewPudge(cfg.DBFolder.DBPath)
	if err != nil {
		logger.Error.Fatalf("failed to create pudge object: %v", err)
	}
	route := routes.NewRoutes(pudge, logger, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		fmt.Print("Select operation:\n 1 - add new vacancy request\n 2 - get stats\n")
		var operation int
		fmt.Fscan(os.Stdin, &operation)
		
		switch operation {
		case 1:
			route.Create(ctx)
		case 2:
			route.GetStats(ctx)
		default:
			fmt.Printf("Invalid operation:%d\n", operation)
		}
	}
}

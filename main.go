package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/routes"
	"github.com/artnikel/vacancystats/internal/storage"
	"github.com/artnikel/vacancystats/internal/utils"
	"github.com/recoilme/pudge"
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

	store, err := storage.NewPudge(cfg.DBFolder.DBPath)
	if err != nil {
		logger.Error.Fatalf("failed to create pudge object: %v", err)
	}
	route := routes.NewRoutes(store, logger, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("\nshutting down ...")
		cancel()
		if err := pudge.CloseAll(); err != nil {
			logger.Error.Fatalf("pudge shutdown error: %v", err)
		}
		os.Exit(0)
	}()

	for {
		fmt.Print("\nselect operation:\n 1 - add new vacancy response\n 2 - get stats\n 3 - delete response\n 4 - update status\n")
		var operation int
		fmt.Fscan(os.Stdin, &operation)

		select {
		case <-ctx.Done():
			return
		default:
			switch operation {
			case 1:
				utils.ClearConsole()
				route.Create(ctx)
			case 2:
				utils.ClearConsole()
				route.GetStats(ctx)
			case 3:
				utils.ClearConsole()
				route.Delete(ctx)
			case 4:
				utils.ClearConsole()
				route.UpdateStatus(ctx)
			default:
				utils.ClearConsole()
				fmt.Printf("\ninvalid operation:%d\n", operation)
			}
		}
	}
}

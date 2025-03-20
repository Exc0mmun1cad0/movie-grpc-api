package main

import (
	"context"
	"fmt"
	"log/slog"
	"movie-service/internal/app"
	"movie-service/internal/config"
	"movie-service/pkg/postgres"
	"movie-service/pkg/sl"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := app.SetupLogger(cfg.MovieService.Env)

	log.Debug("example")

	log.Info(
		"starting app",
		slog.String("env", cfg.MovieService.Env),
		slog.String("address", fmt.Sprintf("localhost:%d", cfg.MovieService.GRPCPort)),
	)

	log.Info("initializing connection to postgres")
	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to db(postgres)", sl.Err(err))
		os.Exit(1)
	}
	log.Info("connected to postgres")

	log.Info("initializing app")
	application := app.New(ctx, log, cfg.MovieService, db)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log.Info("started app")
	go application.GRPCServer.MustRun()

	<-stop
	log.Info("stopping app")

	application.GRPCServer.Stop()

	postgres.MustClose(db)

	log.Info("app stopped")
}

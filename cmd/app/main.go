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
	application, err := app.New(ctx, log, cfg.MovieService, db)
	if err != nil {
		log.Error("failed to init app", sl.Err(err))
		os.Exit(1)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log.Info("started app")
	go application.GRPCServer.MustRun()
	go application.GRPCGateway.MustRun()

	<-stop
	log.Info("stopping app")

	application.GRPCServer.Stop()
	log.Info("grpc server stopped")

	application.GRPCGateway.Stop()
	log.Info("grpc gateway stopped")

	postgres.MustClose(db)
	log.Info("closed connection to postgres db")

	log.Info("stopped app")
}

package app

import (
	"context"
	"fmt"
	"log/slog"
	grpcapp "movie-service/internal/app/grpc"
	"movie-service/internal/app/grpcgateway"
	"movie-service/internal/config"
	repo "movie-service/internal/repository/postgres"
	"movie-service/internal/service/movieservice"

	"github.com/jmoiron/sqlx"
)

type App struct {
	ctx         context.Context
	log         *slog.Logger
	GRPCServer  *grpcapp.App
	GRPCGateway *grpcgateway.Gateway
	DB          *sqlx.DB
}

func New(ctx context.Context, log *slog.Logger, cfg config.MovieService, db *sqlx.DB) (*App, error) {
	const op = "app.New"

	movieRepo := repo.New(db)
	movieService := movieservice.New(movieRepo)

	grpcApp := grpcapp.New(ctx, log, movieService, cfg.GRPCPort)
	grpcGateway, err := grpcgateway.New(ctx, log, cfg.HTTPPort, cfg.GRPCPort)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create grpc gateway server: %w", op, err)
	}

	return &App{
		ctx:         ctx,
		log:         log,
		GRPCServer:  grpcApp,
		GRPCGateway: grpcGateway,
		DB:          db,
	}, nil
}

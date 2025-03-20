package app

import (
	"context"
	"log/slog"
	grpcapp "movie-service/internal/app/grpc"
	"movie-service/internal/config"
	repo "movie-service/internal/repository/postgres"
	"movie-service/internal/service/movieservice"

	"github.com/jmoiron/sqlx"
)

type App struct {
	ctx        context.Context
	log        *slog.Logger
	GRPCServer *grpcapp.App
	DB         *sqlx.DB
}

func New(ctx context.Context, log *slog.Logger, cfg config.MovieService, db *sqlx.DB) *App {
	movieRepo := repo.New(db)
	movieService := movieservice.New(movieRepo)

	grpcApp := grpcapp.New(ctx, log, movieService, cfg.GRPCPort)

	return &App{
		ctx:        ctx,
		log:        log,
		GRPCServer: grpcApp,
		DB:         db,
	}
}

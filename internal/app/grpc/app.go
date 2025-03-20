package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	moviegrpc "movie-service/internal/transport/grpc"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	ctx        context.Context
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       uint16
}

func New(ctx context.Context, log *slog.Logger, movieService moviegrpc.Service, port uint16) *App {
	gRPCServer := grpc.NewServer() // TODO: add here interceptors

	moviegrpc.Register(gRPCServer, log, movieService)

	// TODO: add healthcheck

	return &App{
		ctx:        ctx,
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {

	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: failed to create listener: %w", op, err)
	}

	a.log.Info(fmt.Sprintf("grpc server started on port: %d", a.port), slog.String("op", op))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: failed to run grpc server: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.Info("stopping grpc server...", slog.String("op", op))

	a.gRPCServer.GracefulStop()
}

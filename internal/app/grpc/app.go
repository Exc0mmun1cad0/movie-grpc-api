package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	moviegrpc "movie-service/internal/transport/grpc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	ctx        context.Context
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       uint16
}

func New(ctx context.Context, log *slog.Logger, movieService moviegrpc.Service, port uint16) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			func(
				ctx context.Context,
				req any,
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler,
			) (resp any, err error) {
				return moviegrpc.LoggingInterceptor(log)(ctx, req, info, handler)
			}),
	)

	// Register movie service
	moviegrpc.Register(gRPCServer, log, movieService)

	// TODO: add healthcheck
	// Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthServer)
	healthServer.SetServingStatus("movie-service", grpc_health_v1.HealthCheckResponse_SERVING)

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

package grpcgateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"movie-service/pkg/pb"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	ctx        context.Context
	httpServer *http.Server
	log        *slog.Logger

	httpPort uint16
	grpcPort uint16
}

func New(ctx context.Context, log *slog.Logger, httpPort, grpcPort uint16) (*Gateway, error) {
	const op = "grpcgateway.New"

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterMovieServiceHandlerFromEndpoint(
		ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to register handler for grpc gateway: %w", op, err)
	}

	return &Gateway{
		ctx: ctx,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", httpPort),
			Handler: mux,
		},
		log:      log,
		grpcPort: grpcPort,
		httpPort: httpPort,
	}, nil
}

func (g *Gateway) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

func (g *Gateway) Run() error {
	const op = "grpcgateway.Run"

	g.log.Info(
		fmt.Sprintf("http server (grpc gateway) started on port: %d", g.httpPort),
		slog.String("op", op),
	)

	if err := g.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: failed to run http (grpc gateway) server: %w", op, err)
	}

	return nil
}

func (g *Gateway) Stop() {
	const op = "grpcgateway.Stop"

	g.log.Info(
		"stopping http server(grpc gateway)...",
		slog.String("op", op),
	)

	if err := g.httpServer.Shutdown(g.ctx); err != nil {
		panic(err)
	}
}

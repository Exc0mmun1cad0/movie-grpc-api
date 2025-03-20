package moviegrpc

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type key string

const (
	reqIDKey key = "request_id"
)

func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Generate request ID
		requestID := uuid.New().String()
		log := log.With(slog.String(string(reqIDKey), requestID))

		log.Info(
			"Got new request",
			slog.String("Method", info.FullMethod),
			slog.Any("Body", req),
		)

		// Pass it to handler through the context
		ctx = context.WithValue(ctx, reqIDKey, requestID)

		m, err := handler(ctx, req)

		log.Info(
			"Handled request",
			slog.Any("Response", m),
		)

		return m, err
	}
}

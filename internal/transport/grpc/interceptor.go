package moviegrpc

import (
	"context"
	"log/slog"
	"movie-service/pkg/sl"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type key string

const (
	reqIDKey key = "request_id"
)

func LoggingUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
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
			"Got new unary request",
			slog.String("Method", info.FullMethod),
			slog.Any("Body", req),
		)

		// Pass it to handler through the context
		ctx = context.WithValue(ctx, reqIDKey, requestID)
		var m any
		start := time.Now()

		defer func() {
			log.Info(
				"Request completed",
				slog.Any("Response", m),
				slog.String("duration", time.Since(start).String()),
			)
		}()

		m, err := handler(ctx, req)
		if err != nil {
			log.Error("Failed to handle request", sl.Err(err))
		}

		return m, err
	}
}

func LoggingStreamInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Generate request ID
		requestID := uuid.New().String()
		log := log.With(slog.String(string(reqIDKey), requestID))

		log.Info(
			"Got new stream request",
			slog.String("Method", info.FullMethod),
			slog.Any("Body", srv),
		)

		// Create new context with request ID inside
		ctx := ss.Context()
		ctx = context.WithValue(ctx, reqIDKey, requestID)

		wrappedStream := &serverStreamWrapper{
			ServerStream: ss,
			ctx:          ctx,
		}

		start := time.Now()

		defer func() {
			log.Info(
				"Stream request completed",
				slog.String("duration", time.Since(start).String()),
			)
		}()

		err := handler(srv, wrappedStream)
		if err != nil {
			log.Error("Failed to sream response", sl.Err(err))
		} else {
			log.Info("Finished stream")
		}

		return nil
	}
}

type serverStreamWrapper struct {
	ctx context.Context
	grpc.ServerStream
}

func (w *serverStreamWrapper) Context() context.Context {
	return w.ctx
}

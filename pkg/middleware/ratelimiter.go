package middleware

import (
	"context"
	"log"

	"github.com/sunnyjain123/go-ratelimiter/internal/ratelimiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RateLimiterUnaryInterceptor(svc *ratelimiter.Service) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// For now, static key (later: metadata / auth / IP)
		key := "global"

		decision, err := svc.Check(ctx, key, info.FullMethod)
		if err != nil {
			return nil, status.Error(codes.Internal, "rate limiter error")
		}

		if decision == ratelimiter.Deny {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		log.Printf("rate-limiter allowed: %s", info.FullMethod)

		return handler(ctx, req)
	}
}

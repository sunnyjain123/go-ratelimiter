package ratelimiter

import (
	"context"

	ratelimiterv1 "github.com/sunnyjain123/go-ratelimiter/client/proto/ratelimiter/v1"
)

type RateLimiterService struct {
	ratelimiterv1.RateLimiterServiceServer
}

func NewRateLimiterService() *RateLimiterService {
	return &RateLimiterService{}
}

func (s *RateLimiterService) Ping(
	ctx context.Context,
	req *ratelimiterv1.PingRequest,
) (*ratelimiterv1.PingResponse, error) {

	return &ratelimiterv1.PingResponse{
		Reply: "pong: " + req.Message,
	}, nil
}

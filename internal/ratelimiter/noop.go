package ratelimiter

import "context"

type NoopLimiter struct{}

func NewNoopLimiter() *NoopLimiter {
	return &NoopLimiter{}
}

func (l *NoopLimiter) Allow(ctx context.Context, key string) (Decision, error) {
	return Allow, nil
}

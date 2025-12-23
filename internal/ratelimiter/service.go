package ratelimiter

import "context"

type Service struct {
	limiters map[string]Limiter
	defaultL Limiter
}

func NewService(defaultLimiter Limiter, perMethod map[string]Limiter) *Service {
	return &Service{
		defaultL: defaultLimiter,
		limiters: perMethod,
	}
}

func (s *Service) Check(ctx context.Context, key, method string) (Decision, error) {
	if limiter, ok := s.limiters[method]; ok {
		return limiter.Allow(ctx, key)
	}
	return s.defaultL.Allow(ctx, key)
}

package ratelimiter

import "context"

type Decision int

const (
	Allow Decision = iota
	Deny
)

type Limiter interface {
	Allow(ctx context.Context, key string) (Decision, error)
}

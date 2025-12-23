package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens         float64
	lastRefillTime time.Time
}

type TokenBucketLimiter struct {
	capacity   float64
	refillRate float64 // tokens per second
	mu         sync.Mutex
	buckets    map[string]*tokenBucket
}

func NewTokenBucketLimiter(capacity int, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   float64(capacity),
		refillRate: refillRate,
		buckets:    make(map[string]*tokenBucket),
	}
}

func (l *TokenBucketLimiter) Allow(ctx context.Context, key string) (Decision, error) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.buckets[key]
	if !exists {
		l.buckets[key] = &tokenBucket{
			tokens:         l.capacity - 1,
			lastRefillTime: now,
		}
		return Allow, nil
	}

	elapsed := now.Sub(b.lastRefillTime).Seconds()
	refill := elapsed * l.refillRate

	b.tokens = min(l.capacity, b.tokens+refill)
	b.lastRefillTime = now

	if b.tokens < 1 {
		return Deny, nil
	}

	b.tokens--
	return Allow, nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

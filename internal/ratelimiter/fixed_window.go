package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type fixedWindowEntry struct {
	count     int
	expiresAt time.Time
}

type FixedWindowLimiter struct {
	limit  int
	window time.Duration

	mu     sync.Mutex
	bucket map[string]*fixedWindowEntry
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:  limit,
		window: window,
		bucket: make(map[string]*fixedWindowEntry),
	}
}

func (l *FixedWindowLimiter) Allow(ctx context.Context, key string) (Decision, error) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	entry, exists := l.bucket[key]

	if !exists || now.After(entry.expiresAt) {
		l.bucket[key] = &fixedWindowEntry{
			count:     1,
			expiresAt: now.Add(l.window),
		}
		return Allow, nil
	}

	if entry.count >= l.limit {
		return Deny, nil
	}

	entry.count++
	return Allow, nil
}

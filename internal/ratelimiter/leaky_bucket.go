package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type leakyBucket struct {
	level      int
	lastLeaked time.Time
}

type LeakyBucketLimiter struct {
	capacity int
	leakRate int // requests per second
	mu       sync.Mutex
	buckets  map[string]*leakyBucket
}

func NewLeakyBucketLimiter(capacity, leakRate int) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		capacity: capacity,
		leakRate: leakRate,
		buckets:  make(map[string]*leakyBucket),
	}
}

func (l *LeakyBucketLimiter) Allow(ctx context.Context, key string) (Decision, error) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.buckets[key]
	if !exists {
		l.buckets[key] = &leakyBucket{
			level:      1,
			lastLeaked: now,
		}
		return Allow, nil
	}

	elapsed := now.Sub(b.lastLeaked).Seconds()
	leaked := int(elapsed * float64(l.leakRate))

	if leaked > 0 {
		b.level -= leaked
		if b.level < 0 {
			b.level = 0
		}
		b.lastLeaked = now
	}

	if b.level >= l.capacity {
		return Deny, nil
	}

	b.level++
	return Allow, nil
}

package rate

import (
	"context"
	"sync"
	"time"
)

// RateLimiter defines an interface for rate limiting strategies.
type RateLimiter interface {
	Allow(ctx context.Context, bucket string, limit int, window time.Duration) (allowed bool, remaining int, resetAt time.Time, err error)
}

// InMemorySlidingWindow is a simple per-process sliding-window limiter.
type InMemorySlidingWindow struct {
	mu   sync.Mutex
	hits map[string][]time.Time
}

// NewInMemory creates a new in-memory rate limiter.
func NewInMemory() *InMemorySlidingWindow {
	return &InMemorySlidingWindow{hits: make(map[string][]time.Time)}
}

// Allow implements RateLimiter using a timestamp slice per bucket.
func (l *InMemorySlidingWindow) Allow(ctx context.Context, bucket string, limit int, window time.Duration) (bool, int, time.Time, error) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	arr := l.hits[bucket]
	cutoff := now.Add(-window)
	// drop old entries
	j := 0
	for ; j < len(arr); j++ {
		if arr[j].After(cutoff) {
			break
		}
	}
	if j > 0 {
		arr = arr[j:]
	} else if len(arr) > 0 && arr[0].Before(cutoff) {
		arr = nil
	}

	if len(arr) >= limit {
		reset := arr[0].Add(window)
		l.hits[bucket] = arr
		return false, 0, reset, nil
	}

	arr = append(arr, now)
	l.hits[bucket] = arr
	remaining := limit - len(arr)
	reset := arr[0].Add(window)
	return true, remaining, reset, nil
}

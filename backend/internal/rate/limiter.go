package rate

import (
	"context"
	"fmt"
	"sort"
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
	if err := ctx.Err(); err != nil {
		return false, 0, time.Time{}, err
	}
	if limit <= 0 || window <= 0 {
		return false, 0, time.Time{}, fmt.Errorf("invalid params: limit=%d window=%s", limit, window)
	}

	now := time.Now()
	cutoff := now.Add(-window)

	l.mu.Lock()
	defer l.mu.Unlock()
	arr := l.hits[bucket]

	// binary search for first valid
	idx := sort.Search(len(arr), func(i int) bool { return arr[i].After(cutoff) })
	if idx > 0 {
		arr = arr[idx:]
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

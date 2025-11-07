package rate

import (
	"context"
	"testing"
	"time"
)

func TestInMemorySlidingWindow(t *testing.T) {
	limiter := NewInMemory()
	ctx := context.Background()
	bucket := "test:bucket"
	limit := 3
	window := 200 * time.Millisecond

	// First 3 should pass
	for i := 0; i < limit; i++ {
		allowed, remaining, _, err := limiter.Allow(ctx, bucket, limit, window)
		if err != nil || !allowed {
			t.Fatalf("expected allowed on iteration %d got err=%v allowed=%v", i, err, allowed)
		}
		if remaining != (limit - (i + 1)) {
			t.Errorf("remaining mismatch: got %d want %d", remaining, limit-(i+1))
		}
	}
	// 4th should block
	allowed, remaining, _, err := limiter.Allow(ctx, bucket, limit, window)
	if err != nil || allowed || remaining != 0 {
		t.Errorf("expected deny on 4th call: allowed=%v remaining=%d err=%v", allowed, remaining, err)
	}
	// Wait for window to expire
	time.Sleep(window + 20*time.Millisecond)
	allowed, remaining, _, err = limiter.Allow(ctx, bucket, limit, window)
	if err != nil || !allowed || remaining != limit-1 {
		t.Errorf("expected allow after window reset: allowed=%v remaining=%d err=%v", allowed, remaining, err)
	}
}

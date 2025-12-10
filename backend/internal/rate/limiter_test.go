package rate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInMemorySlidingWindow(t *testing.T) {
	t.Parallel()

	limiter := NewInMemory()
	ctx := context.Background()
	bucket := "test:bucket"
	limit := 3
	window := 200 * time.Millisecond

	for i := 0; i < limit; i++ {
		allowed, remaining, _, err := limiter.Allow(ctx, bucket, limit, window)
		require.NoErrorf(t, err, "iteration %d should not error", i)
		require.Truef(t, allowed, "iteration %d should be allowed", i)
		require.Equalf(t, limit-(i+1), remaining, "iteration %d remaining mismatch", i)
	}

	allowed, remaining, _, err := limiter.Allow(ctx, bucket, limit, window)
	require.NoError(t, err)
	require.False(t, allowed)
	require.Zero(t, remaining)

	time.Sleep(window + 20*time.Millisecond)
	allowed, remaining, _, err = limiter.Allow(ctx, bucket, limit, window)
	require.NoError(t, err)
	require.True(t, allowed)
	require.Equal(t, limit-1, remaining)
}

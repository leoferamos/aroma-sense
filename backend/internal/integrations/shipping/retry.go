package shipping

import (
	"context"
	"errors"
	"net"
)

// shouldRetryError indicates whether a request error is transient and worth retrying.
func shouldRetryError(err error) bool {
	// context deadline
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// net timeout
	var ne net.Error
	if errors.As(err, &ne) {
		if ne.Timeout() {
			return true
		}
	}
	return false
}

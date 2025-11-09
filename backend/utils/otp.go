package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP generates a cryptographically secure 6-digit numeric code.
func GenerateOTP() (string, error) {
	// Generate a random number between 0 and 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}

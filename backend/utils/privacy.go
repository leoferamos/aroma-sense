package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// MaskEmail masks an email address for privacy (LGPD compliance)
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local := parts[0]
	domain := parts[1]

	// Mask local part (show first and last character, mask middle)
	if len(local) > 2 {
		local = string(local[0]) + strings.Repeat("*", len(local)-2) + string(local[len(local)-1])
	} else if len(local) == 2 {
		local = string(local[0]) + "*"
	}

	// Mask domain part (show first and last character, mask middle)
	if dotIndex := strings.LastIndex(domain, "."); dotIndex > 0 {
		domainPart := domain[:dotIndex]
		tld := domain[dotIndex:]

		if len(domainPart) > 2 {
			domainPart = string(domainPart[0]) + strings.Repeat("*", len(domainPart)-2) + string(domainPart[len(domainPart)-1])
		} else if len(domainPart) == 2 {
			domainPart = string(domainPart[0]) + "*"
		}

		domain = domainPart + tld
	}

	return local + "@" + domain
}

// HashEmailForLogging creates a SHA256 hash of email for secure logging (LGPD compliance)
func HashEmailForLogging(email string) string {
	if email == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	return fmt.Sprintf("%x", hash)[:16] // First 16 chars of hash
}

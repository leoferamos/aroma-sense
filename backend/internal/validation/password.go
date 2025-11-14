package validation

import (
	"errors"
	"regexp"
	"strings"
)

var (
	reLower = regexp.MustCompile(`[a-z]`)
	reUpper = regexp.MustCompile(`[A-Z]`)
	reDigit = regexp.MustCompile(`[0-9]`)
)

const MinPasswordLen = 8

// ValidatePassword applies complexity rules
func ValidatePassword(pw, email string) error {
	pw = strings.TrimSpace(pw)
	if len(pw) < MinPasswordLen {
		return errors.New("password must be at least 8 characters long")
	}
	if !reLower.MatchString(pw) || !reUpper.MatchString(pw) || !reDigit.MatchString(pw) {
		return errors.New("password must contain at least one lowercase letter, one uppercase letter, and one digit")
	}
	return nil
}

package ai

import (
	"regexp"
	"strings"
)

// SanitizeUserMessage removes PII and truncates overly long inputs to keep prompts safe and on-topic.
func SanitizeUserMessage(input string, maxRunes int) string {
	s := strings.TrimSpace(input)
	if maxRunes <= 0 {
		maxRunes = 800
	}
	emailRe := regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	phoneRe := regexp.MustCompile(`(?i)\+?\d[\d\s().\-]{7,}\d`)
	urlRe := regexp.MustCompile(`(?i)\bhttps?://[^\s]+`)

	s = emailRe.ReplaceAllString(s, "[redacted-email]")
	s = phoneRe.ReplaceAllString(s, "[redacted-phone]")
	s = urlRe.ReplaceAllString(s, "[link]")

	// Truncate by runes to avoid breaking multi-byte chars
	return truncateRunes(s, maxRunes)
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	// Fast path
	if len([]rune(s)) <= max {
		return s
	}
	r := []rune(s)
	return string(r[:max])
}

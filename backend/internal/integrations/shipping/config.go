package shipping

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config centralizes shipping provider configuration.
type Config struct {
	BaseURL     string
	QuotesPath  string
	StaticToken string
	UserAgent   string
	Services    string
	OriginCEP   string
	Timeout     time.Duration
}

// LoadShippingConfigFromEnv loads configuration for the provider from environment variables.
func LoadShippingConfigFromEnv() (Config, error) {
	cfg := Config{}

	// Base URL and paths
	cfg.BaseURL = withDefault(getenvFirstNonEmpty("SHIPPING_BASE_URL", "SUPERFRETE_BASE_URL"), "https://sandbox.superfrete.com")
	cfg.QuotesPath = withDefault(getenvFirstNonEmpty("SHIPPING_QUOTES_PATH", "SUPERFRETE_QUOTES_PATH"), "/api/v0/calculator")

	// Client identification and services
	cfg.UserAgent = getenvFirstNonEmpty("SHIPPING_USER_AGENT", "SUPERFRETE_USER_AGENT")
	cfg.Services = withDefault(getenvFirstNonEmpty("SHIPPING_SERVICES", "SUPERFRETE_SERVICES"), "1,2,17")
	cfg.OriginCEP = os.Getenv("SHIPPING_ORIGIN_CEP")

	// Token (trim and optional removal of Bearer prefix)
	cfg.StaticToken = readBearerToken("SHIPPING_TOKEN", "SUPERFRETE_TOKEN")

	// Timeout in seconds
	cfg.Timeout = readTimeoutSecondsFirstNonEmpty("SHIPPING_TIMEOUT", "SUPERFRETE_TIMEOUT")
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	if cfg.BaseURL == "" || cfg.OriginCEP == "" || cfg.StaticToken == "" || cfg.UserAgent == "" {
		return cfg, errors.New("shipping provider env not fully configured: base url, origin CEP, token and user agent required")
	}
	if strings.HasPrefix(strings.TrimSpace(cfg.StaticToken), "REPLACE_WITH") {
		return cfg, errors.New("shipping provider token placeholder detected; set the token environment variable with a valid value")
	}
	return cfg, nil
}

// getenvFirstNonEmpty returns the first non-empty value among the provided keys.
func getenvFirstNonEmpty(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// withDefault returns def if s is empty, otherwise s.
func withDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// readBearerToken reads the first non-empty token from the provided env keys and sanitizes it.
func readBearerToken(keys ...string) string {
	tok := getenvFirstNonEmpty(keys...)
	if tok == "" {
		return ""
	}
	tok = strings.Trim(tok, " \t\n\r\"'")
	low := strings.ToLower(tok)
	if strings.HasPrefix(low, "bearer ") {
		tok = strings.TrimSpace(tok[7:])
	}
	return tok
}

// readTimeoutSecondsFirstNonEmpty parses the first non-empty env var (seconds) into a duration.
func readTimeoutSecondsFirstNonEmpty(keys ...string) time.Duration {
	for _, k := range keys {
		if s := os.Getenv(k); s != "" {
			if n, err := strconv.Atoi(s); err == nil && n > 0 {
				return time.Duration(n) * time.Second
			}
		}
	}
	return 0
}

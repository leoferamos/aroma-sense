package shipping

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Config centralizes all shipping-related configuration.
type Config struct {
	BaseURL      string
	TokenURL     string
	ClientID     string
	ClientSecret string
	QuotesPath   string        // optional; defaults to "/quotes"
	OriginCEP    string        // required for calculating quotes
	Timeout      time.Duration // optional; default 15s
}

// LoadShippingConfigFromEnv loads configuration from environment variables.
func LoadShippingConfigFromEnv() (Config, error) {
	cfg := Config{
		BaseURL:      os.Getenv("SHIPPING_BASE_URL"),
		TokenURL:     os.Getenv("SHIPPING_TOKEN_URL"),
		ClientID:     os.Getenv("SHIPPING_CLIENT_ID"),
		ClientSecret: os.Getenv("SHIPPING_CLIENT_SECRET"),
		QuotesPath:   os.Getenv("SHIPPING_QUOTES_PATH"),
		OriginCEP:    os.Getenv("SHIPPING_ORIGIN_CEP"),
	}
	if cfg.QuotesPath == "" {
		cfg.QuotesPath = "/quotes"
	}
	if ts := os.Getenv("SHIPPING_TIMEOUT"); ts != "" {
		if n, err := strconv.Atoi(ts); err == nil && n > 0 {
			cfg.Timeout = time.Duration(n) * time.Second
		}
	}
	if cfg.BaseURL == "" || cfg.TokenURL == "" || cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.OriginCEP == "" {
		return cfg, errors.New("shipping env not fully configured")
	}
	return cfg, nil
}

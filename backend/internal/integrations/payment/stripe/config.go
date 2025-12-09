package stripe

import (
	"fmt"
	"os"
)

// Config holds Stripe credentials.
type Config struct {
	SecretKey     string
	WebhookSecret string
}

// LoadConfigFromEnv reads required Stripe variables from environment.
func LoadConfigFromEnv() (*Config, error) {
	secret := os.Getenv("STRIPE_SECRET_KEY")
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	if secret == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY not set")
	}

	return &Config{
		SecretKey:     secret,
		WebhookSecret: webhookSecret,
	}, nil
}

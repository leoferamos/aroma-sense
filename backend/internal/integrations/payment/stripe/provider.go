package stripe

import (
	"context"
	"encoding/json"
	"fmt"

	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"

	"github.com/leoferamos/aroma-sense/internal/service"
)

// Provider implements a minimal Stripe payment provider for PaymentIntent creation and webhook parsing.
type Provider struct {
	webhookSecret string
}

// NewProvider returns a configured Stripe Provider.
func NewProvider(cfg *Config) *Provider {
	stripe.Key = cfg.SecretKey
	return &Provider{webhookSecret: cfg.WebhookSecret}
}

// PaymentIntentResult contains the minimal fields returned to callers.
// CreatePaymentIntent creates a payment intent for the given amount and currency.
func (p *Provider) CreatePaymentIntent(ctx context.Context, params service.PaymentIntentParams) (*service.PaymentIntentResult, error) {
	piParams := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(params.Amount),
		Currency:           stripe.String(params.Currency),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Metadata:           params.Metadata,
	}

	if params.CustomerEmail != "" {
		piParams.ReceiptEmail = stripe.String(params.CustomerEmail)
	}

	intent, err := paymentintent.New(piParams)
	if err != nil {
		return nil, fmt.Errorf("stripe create payment intent: %w", err)
	}
	return &service.PaymentIntentResult{ID: intent.ID, ClientSecret: intent.ClientSecret}, nil
}

// ParseWebhook validates signature and returns a normalized payload.
func (p *Provider) ParseWebhook(payload []byte, signature string) (*service.PaymentWebhookPayload, error) {
	if p.webhookSecret == "" {
		return nil, fmt.Errorf("stripe webhook secret not configured")
	}
	event, err := webhook.ConstructEvent(payload, signature, p.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("stripe webhook validation failed: %w", err)
	}

	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed", "payment_intent.canceled", "payment_intent.processing":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return nil, fmt.Errorf("stripe webhook unmarshal: %w", err)
		}
		status := string(pi.Status)
		return &service.PaymentWebhookPayload{
			IntentID:      pi.ID,
			Status:        status,
			Amount:        pi.Amount,
			Currency:      string(pi.Currency),
			CustomerEmail: pi.ReceiptEmail,
			Metadata:      pi.Metadata,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported webhook event: %s", event.Type)
	}
}

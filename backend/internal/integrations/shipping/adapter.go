package shipping

import (
	"context"
	"strings"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/utils"
)

// Provider implements a ShippingProvider
type Provider struct {
	client        *Client
	quotesPath    string
	token         string
	userAgent     string
	services      string
	cache         QuoteCache
	retryAttempts int
	retryBackoff  time.Duration
}

func NewProvider(client *Client) *Provider {
	return &Provider{
		client:        client,
		quotesPath:    "/quotes",
		cache:         NewInMemoryQuoteCache(60 * time.Second),
		retryAttempts: 2, // 1 initial + 1 retry
		retryBackoff:  300 * time.Millisecond,
	}
}

// WithQuotesPath allows overriding the quotes path.
func (p *Provider) WithQuotesPath(path string) *Provider {
	if path != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		p.quotesPath = path
	}
	return p
}

// WithStaticAuth sets the Bearer token and User-Agent required by the shipping provider.
func (p *Provider) WithStaticAuth(token, userAgent string) *Provider {
	p.token = token
	p.userAgent = userAgent
	return p
}

// WithServices sets the provider services string.
func (p *Provider) WithServices(services string) *Provider {
	if services != "" {
		p.services = services
	}
	return p
}

// WithCacheTTL allows tuning the in-memory cache TTL for quotes.
func (p *Provider) WithCacheTTL(d time.Duration) *Provider {
	if d > 0 {
		if c, ok := p.cache.(*InMemoryQuoteCache); ok {
			c.SetTTL(d)
		}
	}
	return p
}

// WithRetry allows configuring retry attempts and backoff duration.
func (p *Provider) WithRetry(attempts int, backoff time.Duration) *Provider {
	if attempts >= 1 {
		p.retryAttempts = attempts
	}
	if backoff > 0 {
		p.retryBackoff = backoff
	}
	return p
}

// GetQuotes retrieves shipping quotes from the provider.
func (p *Provider) GetQuotes(ctx context.Context, userID string, originPostalCode, destPostalCode string, parcels []model.Parcel, insuredValue float64) ([]dto.ShippingOption, error) {
	// Build payload
	reqBody := p.buildQuoteRequest(originPostalCode, destPostalCode, parcels, insuredValue)
	key := p.cacheKey(reqBody)
	if opts, ok := p.cache.Get(key); ok {
		return opts, nil
	}
	// Perform HTTP with retry
	items, err := p.fetchQuotes(ctx, reqBody)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		p.cache.Set(key, []dto.ShippingOption{})
		return []dto.ShippingOption{}, nil
	}
	opts := mapProviderQuotes(items)
	p.cache.Set(key, opts)
	return opts, nil
}

// cacheKey builds a stable key for the given request body.
func (p *Provider) cacheKey(r quoteRequest) string {
	// Round floats to reduce key cardinality
	w := utils.FormatFloatTrim(r.Package.Weight, 3)
	h := utils.FormatFloatTrim(r.Package.Height, 1)
	wi := utils.FormatFloatTrim(r.Package.Width, 1)
	l := utils.FormatFloatTrim(r.Package.Length, 1)
	iv := utils.FormatFloatTrim(r.Options.InsuranceValue, 2)
	return r.From.PostalCode + "|" + r.To.PostalCode + "|" + w + "|" + h + "|" + wi + "|" + l + "|" + iv + "|" + p.services
}

package shipping

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// fetchQuotes performs the HTTP call with a small retry policy and returns provider quotes.
func (p *Provider) fetchQuotes(ctx context.Context, body quoteRequest) ([]providerQuote, error) {
	attempts := p.retryAttempts
	if attempts < 1 {
		attempts = 1
	}
	var lastErr error
	for try := 1; try <= attempts; try++ {
		req, err := p.client.NewJSONRequest(http.MethodPost, p.quotesPath, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+p.token)
		req.Header.Set("User-Agent", p.userAgent)
		req = req.WithContext(ctx)

		resp, err := p.client.Do(req)
		if err != nil {
			if shouldRetryError(err) && try < attempts {
				time.Sleep(p.retryBackoff)
				continue
			}
			return nil, err
		}
		// Handle response
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			_ = resp.Body.Close()
			msg := strings.TrimSpace(string(bodyBytes))
			if msg == "" {
				msg = resp.Status
			} else {
				msg = resp.Status + ": " + msg
			}
			lastErr = fmt.Errorf("shipping quotes failed: %s", msg)
			if (resp.StatusCode >= 500 || resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode == http.StatusTooManyRequests) && try < attempts {
				time.Sleep(p.retryBackoff)
				continue
			}
			// non-retryable
			return nil, lastErr
		}
		var items []providerQuote
		decErr := json.NewDecoder(resp.Body).Decode(&items)
		_ = resp.Body.Close()
		if decErr != nil {
			lastErr = decErr
			if try < attempts {
				time.Sleep(p.retryBackoff)
				continue
			}
			return nil, lastErr
		}
		return items, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("unexpected state in fetchQuotes")
}

package shipping

import (
	"os"
)

// NewProviderFromEnv builds a Provider using environment variables.
func NewProviderFromEnv() *Provider {
	cli, err := NewClientFromEnv()
	if err != nil {
		return nil
	}
	p := NewProvider(cli)
	if qp := os.Getenv("SHIPPING_QUOTES_PATH"); qp != "" {
		p.WithQuotesPath(qp)
	}
	return p
}

// OriginCEPFromEnv returns the origin CEP from environment variables.
func OriginCEPFromEnv() string {
	return os.Getenv("SHIPPING_ORIGIN_CEP")
}

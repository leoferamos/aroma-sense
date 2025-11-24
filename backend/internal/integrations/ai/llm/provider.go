package llm

import "context"

// Provider is a minimal interface for generating chat completions.
type Provider interface {
	Generate(ctx context.Context, prompt string, maxTokens int) (string, error)
}

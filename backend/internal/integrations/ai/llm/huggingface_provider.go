package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/leoferamos/aroma-sense/internal/integrations/ai/config"
)

// HuggingFaceProvider implements Provider using Hugging Face Inference API.
type HuggingFaceProvider struct {
	apiKey  string
	model   string
	client  *http.Client
	timeout time.Duration
}

// NewHuggingFaceProvider creates a Provider backed by Hugging Face Inference API.
func NewHuggingFaceProvider(cfg config.Config) Provider {
	return &HuggingFaceProvider{
		apiKey:  cfg.APIKey,
		model:   cfg.LLMModel,
		client:  &http.Client{Timeout: cfg.Timeout},
		timeout: cfg.Timeout,
	}
}

// Generate generates a completion for the given prompt using Hugging Face Inference Providers API.
func (p *HuggingFaceProvider) Generate(ctx context.Context, prompt string, maxTokens int) (string, error) {
	url := "https://router.huggingface.co/v1/chat/completions"

	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  maxTokens,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s, body: %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}

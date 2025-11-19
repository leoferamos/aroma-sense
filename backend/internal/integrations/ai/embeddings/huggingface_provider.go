package embeddings

import (
	"bytes"
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
		model:   cfg.EmbModel,
		client:  &http.Client{Timeout: cfg.Timeout},
		timeout: cfg.Timeout,
	}
}

// Embed generates embeddings for the given texts using Hugging Face API.
func (p *HuggingFaceProvider) Embed(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", p.model)

	payload := map[string]interface{}{
		"inputs": texts,
		"options": map[string]interface{}{
			"wait_for_model": true,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, body: %s", resp.Status, string(body))
	}

	var result []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	embeddings := make([][]float32, len(result))
	for i, item := range result {
		if vec, ok := item.([]interface{}); ok {
			embeddings[i] = make([]float32, len(vec))
			for j, val := range vec {
				if f, ok := val.(float64); ok {
					embeddings[i][j] = float32(f)
				}
			}
		}
	}

	return embeddings, nil
}

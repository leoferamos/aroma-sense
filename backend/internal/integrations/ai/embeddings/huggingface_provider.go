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
	apiKey        string
	model         string
	client        *http.Client
	timeout       time.Duration
	queryPrefix   string
	passagePrefix string
	baseURL       string
}

// NewHuggingFaceProvider creates a Provider backed by Hugging Face Inference API.
func NewHuggingFaceProvider(cfg config.Config) Provider {
	p := &HuggingFaceProvider{
		apiKey:  cfg.APIKey,
		model:   cfg.EmbModel,
		client:  &http.Client{Timeout: cfg.Timeout},
		timeout: cfg.Timeout,
		baseURL: "https://router.huggingface.co/v1",
	}
	// Configure based on model
	p.configureForModel()
	return p
}

// configureForModel sets model-specific configurations
func (p *HuggingFaceProvider) configureForModel() {
	switch p.model {
	case "intfloat/multilingual-e5-large", "intfloat/e5-small-v2":
		p.queryPrefix = "query: "
		p.passagePrefix = "passage: "
		p.baseURL = "https://router.huggingface.co/v1"
	case "BAAI/bge-small-en-v1.5":
		p.queryPrefix = ""
		p.passagePrefix = ""
		p.baseURL = "https://router.huggingface.co/v1"
	case "Qwen/Qwen3-Embedding-8B":
		p.queryPrefix = ""
		p.passagePrefix = ""
		p.baseURL = "https://router.huggingface.co/nebius/v1"
	default:
		p.queryPrefix = ""
		p.passagePrefix = ""
		p.baseURL = "https://router.huggingface.co/v1"
	}
}

// Embed generates embeddings for the given texts using Hugging Face Inference Providers API.
func (p *HuggingFaceProvider) Embed(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	// Apply prefixes if configured
	processedTexts := make([]string, len(texts))
	for i, text := range texts {
		processedTexts[i] = p.passagePrefix + text
	}

	url := p.baseURL + "/embeddings"

	payload := map[string]interface{}{
		"model": p.model,
		"input": processedTexts,
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

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if data, ok := result["data"].([]interface{}); ok {
		embeddings := make([][]float32, len(data))
		for i, item := range data {
			if embeddingData, ok := item.(map[string]interface{}); ok {
				if embedding, ok := embeddingData["embedding"].([]interface{}); ok {
					embeddings[i] = make([]float32, len(embedding))
					for j, val := range embedding {
						if f, ok := val.(float64); ok {
							embeddings[i][j] = float32(f)
						}
					}
				}
			}
		}
		return embeddings, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

// EmbedQuery generates an embedding for a single query text.
func (p *HuggingFaceProvider) EmbedQuery(query string) ([]float32, error) {
	processedQuery := p.queryPrefix + query
	embeddings, err := p.Embed([]string{processedQuery})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return embeddings[0], nil
}

// Configure allows setting model-specific options
func (p *HuggingFaceProvider) Configure(config map[string]interface{}) Provider {
	if queryPrefix, ok := config["queryPrefix"].(string); ok {
		p.queryPrefix = queryPrefix
	}
	if passagePrefix, ok := config["passagePrefix"].(string); ok {
		p.passagePrefix = passagePrefix
	}
	return p
}

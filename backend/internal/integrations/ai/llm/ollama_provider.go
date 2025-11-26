package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OllamaConfig holds minimal configuration for local LLM inference.
type OllamaConfig struct {
	BaseURL   string
	Model     string
	Timeout   time.Duration
	MaxTokens int
}

type ollamaPayload struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Options struct {
		NumPredict int `json:"num_predict,omitempty"`
	} `json:"options"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	// When stream=false Ollama returns a single JSON object with aggregated response
	Done bool `json:"done"`
}

// OllamaProvider implements Provider against a local Ollama server.
type OllamaProvider struct {
	cfg    OllamaConfig
	client *http.Client
}

func NewOllamaProvider(cfg OllamaConfig) *OllamaProvider {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 256
	}
	return &OllamaProvider{cfg: cfg, client: &http.Client{Timeout: cfg.Timeout}}
}

func (p *OllamaProvider) Generate(ctx context.Context, prompt string, maxTokens int) (string, error) {
	if maxTokens <= 0 || maxTokens > p.cfg.MaxTokens {
		maxTokens = p.cfg.MaxTokens
	}
	body := ollamaPayload{Model: p.cfg.Model, Prompt: prompt, Stream: false}
	body.Options.NumPredict = maxTokens
	buf, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/api/generate", p.cfg.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama status %d", resp.StatusCode)
	}
	var out ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Response, nil
}

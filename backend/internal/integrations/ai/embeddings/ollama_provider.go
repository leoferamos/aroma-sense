package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OllamaConfig configures a local Ollama server for embeddings.
type OllamaConfig struct {
	BaseURL string        // e.g., http://localhost:11434
	Model   string        // e.g., "nomic-embed-text" or "all-minilm"
	Timeout time.Duration // default 10s
}

type ollamaProvider struct {
	client *http.Client
	cfg    OllamaConfig
}

// NewOllamaProvider creates a Provider backed by a local Ollama server.
func NewOllamaProvider(cfg OllamaConfig) Provider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:11434"
	}
	if cfg.Model == "" {
		cfg.Model = "nomic-embed-text"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	return &ollamaProvider{client: &http.Client{Timeout: cfg.Timeout}, cfg: cfg}
}

type ollamaEmbReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbResp struct {
	Embedding []float32 `json:"embedding"`
}

// Embed calls Ollama /api/embeddings once per text (simple, reliable on current API).
func (p *ollamaProvider) Embed(texts []string) ([][]float32, error) {
	out := make([][]float32, 0, len(texts))
	for _, t := range texts {
		reqBody, _ := json.Marshal(ollamaEmbReq{Model: p.cfg.Model, Prompt: t})
		url := fmt.Sprintf("%s/api/embeddings", p.cfg.BaseURL)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := p.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("ollama http status %d", resp.StatusCode)
		}
		var r ollamaEmbResp
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return nil, err
		}
		out = append(out, r.Embedding)
	}
	return out, nil
}

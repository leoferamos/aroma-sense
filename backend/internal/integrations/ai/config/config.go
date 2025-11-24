package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Config centralizes AI provider configuration.
type Config struct {
	Provider   string
	LLMBaseURL string
	LLMModel   string
	EmbBaseURL string
	EmbModel   string
	APIKey     string
	Timeout    time.Duration
}

// LoadAIConfigFromEnv loads configuration for AI providers from environment variables.
func LoadAIConfigFromEnv() (Config, error) {
	cfg := Config{}

	// Provider selection
	cfg.Provider = withDefault(getenvFirstNonEmpty("AI_PROVIDER"), "ollama")

	// LLM configuration with backward compatibility
	cfg.LLMBaseURL = withDefault(getenvFirstNonEmpty("AI_LLM_BASE_URL", "OLLAMA_LLM_BASE_URL"), "http://localhost:11434")
	cfg.LLMModel = withDefault(getenvFirstNonEmpty("AI_LLM_MODEL", "OLLAMA_LLM_MODEL"), "tinyllama:latest")

	// Embedding configuration with backward compatibility
	cfg.EmbBaseURL = withDefault(getenvFirstNonEmpty("AI_EMB_BASE_URL", "OLLAMA_EMB_BASE_URL"), cfg.LLMBaseURL) // fallback to LLM base URL
	cfg.EmbModel = withDefault(getenvFirstNonEmpty("AI_EMB_MODEL", "OLLAMA_EMB_MODEL"), "nomic-embed-text:latest")

	// API Key for cloud providers
	cfg.APIKey = getenvFirstNonEmpty("AI_API_KEY", "HUGGINGFACE_API_KEY")

	// Timeout
	cfg.Timeout = readTimeoutSecondsFirstNonEmpty("AI_TIMEOUT")
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	// Validate based on provider
	if cfg.Provider == "huggingface" {
		if cfg.APIKey == "" {
			return cfg, errors.New("AI_API_KEY or HUGGINGFACE_API_KEY required for Hugging Face provider")
		}
	}

	return cfg, nil
}

// readTimeoutSecondsFirstNonEmpty parses the first non-empty env var (seconds) into a duration.
func readTimeoutSecondsFirstNonEmpty(keys ...string) time.Duration {
	for _, k := range keys {
		if s := os.Getenv(k); s != "" {
			if n, err := strconv.Atoi(s); err == nil && n > 0 {
				return time.Duration(n) * time.Second
			}
		}
	}
	return 0
}

// getenvFirstNonEmpty returns the first non-empty value among the provided keys.
func getenvFirstNonEmpty(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// withDefault returns def if s is empty, otherwise s.
func withDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

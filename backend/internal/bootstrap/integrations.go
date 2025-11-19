package bootstrap

import (
	"log"
	"time"

	"github.com/leoferamos/aroma-sense/internal/email"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/config"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/embeddings"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/llm"
	shippingprovider "github.com/leoferamos/aroma-sense/internal/integrations/shipping"
	"github.com/leoferamos/aroma-sense/internal/service"
)

// integrations holds all external integration instances
type integrations struct {
	email    email.EmailService
	shipping *shippingIntegration
	ai       *aiIntegration
}

// shippingIntegration holds shipping-related services
type shippingIntegration struct {
	provider  service.ShippingProvider
	service   service.ShippingService
	originCEP string
}

// aiIntegration holds AI-related providers
type aiIntegration struct {
	llmProvider llm.Provider
	embProvider embeddings.Provider
}

// initializeIntegrations creates all external integration instances
func initializeIntegrations() *integrations {
	return &integrations{
		email:    initializeEmailIntegration(),
		shipping: initializeShippingIntegration(),
		ai:       initializeAIIntegration(),
	}
}

// initializeEmailIntegration initializes email service
func initializeEmailIntegration() email.EmailService {
	smtpConfig := email.LoadSMTPConfigFromEnv()
	if err := smtpConfig.Validate(); err != nil {
		log.Fatalf("SMTP configuration error: %v", err)
	}

	emailService, err := email.NewSMTPEmailService(smtpConfig)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	return emailService
}

// initializeShippingIntegration initializes shipping provider and service
func initializeShippingIntegration() *shippingIntegration {
	cfg, err := shippingprovider.LoadShippingConfigFromEnv()
	if err != nil {
		log.Printf("Shipping configuration not available: %v", err)
		return nil
	}

	cli, err := shippingprovider.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to create shipping client: %v", err)
		return nil
	}

	provider := shippingprovider.NewProvider(cli).
		WithQuotesPath(cfg.QuotesPath).
		WithStaticAuth(cfg.StaticToken, cfg.UserAgent).
		WithServices(cfg.Services)

	return &shippingIntegration{
		provider:  provider,
		originCEP: cfg.OriginCEP,
		// service will be initialized later with repositories
	}
}

// initializeAIIntegration initializes AI providers
func initializeAIIntegration() *aiIntegration {
	cfg, err := config.LoadAIConfigFromEnv()
	if err != nil {
		log.Printf("AI configuration error: %v", err)
		log.Printf("Falling back to default Ollama configuration")
		// Fallback to default Ollama config
		cfg = config.Config{
			Provider:   "ollama",
			LLMBaseURL: "http://localhost:11434",
			LLMModel:   "tinyllama:latest",
			EmbBaseURL: "http://localhost:11434",
			EmbModel:   "nomic-embed-text:latest",
			Timeout:    30 * time.Second,
		}
	}

	var llmProvider llm.Provider
	var embProvider embeddings.Provider

	if cfg.Provider == "huggingface" {
		llmProvider = llm.NewHuggingFaceProvider(cfg)
		embProvider = embeddings.NewHuggingFaceProvider(cfg)
	} else {
		// Default to Ollama
		llmProvider = llm.NewOllamaProvider(llm.OllamaConfig{
			BaseURL: cfg.LLMBaseURL,
			Model:   cfg.LLMModel,
		})
		embProvider = embeddings.NewOllamaProvider(embeddings.OllamaConfig{
			BaseURL: cfg.EmbBaseURL,
			Model:   cfg.EmbModel,
		})
	}

	return &aiIntegration{
		llmProvider: llmProvider,
		embProvider: embProvider,
	}
}

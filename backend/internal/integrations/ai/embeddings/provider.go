package embeddings

type Provider interface {
	Embed(texts []string) ([][]float32, error)
	EmbedQuery(query string) ([]float32, error)
	// Configure allows setting model-specific options
	Configure(config map[string]interface{}) Provider
}

// EmbeddingConfig holds configuration for embedding providers
type EmbeddingConfig struct {
	Model         string
	QueryPrefix   string
	PassagePrefix string
	MaxTokens     int
}

// Noop is a dummy embeddings provider that returns zero-vectors.
type Noop struct {
	Dim int
}

// Embed returns zero-vectors for each input text.
func (n Noop) Embed(texts []string) ([][]float32, error) {
	if n.Dim <= 0 {
		n.Dim = 384
	}
	out := make([][]float32, len(texts))
	for i := range texts {
		out[i] = make([]float32, n.Dim)
	}
	return out, nil
}

// Configure implements Provider interface
func (n Noop) Configure(config map[string]interface{}) Provider {
	if dim, ok := config["dim"].(int); ok {
		n.Dim = dim
	}
	return n
}

// EmbedQuery returns a zero-vector for the query.
func (n Noop) EmbedQuery(query string) ([]float32, error) {
	if n.Dim <= 0 {
		n.Dim = 384
	}
	return make([]float32, n.Dim), nil
}

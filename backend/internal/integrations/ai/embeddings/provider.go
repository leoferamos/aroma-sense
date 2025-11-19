package embeddings

type Provider interface {
	Embed(texts []string) ([][]float32, error)
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

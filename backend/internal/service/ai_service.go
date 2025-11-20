package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/leoferamos/aroma-sense/internal/ai"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// AIService provides low-cost recommendation logic without external calls.
type AIService struct {
	repo  repository.ProductRepository
	mu    sync.Mutex
	cache map[string]cacheEntry
	ttl   time.Duration
}

type cacheEntry struct {
	suggestions []dto.RecommendSuggestion
	expiresAt   time.Time
}

func NewAIService(repo repository.ProductRepository) *AIService {
	return &AIService{
		repo:  repo,
		cache: make(map[string]cacheEntry),
		ttl:   2 * time.Minute,
	}
}

// Recommend performs a full-text retrieval on products using the sanitized message.
func (s *AIService) Recommend(ctx context.Context, rawMessage string, limit int) ([]dto.RecommendSuggestion, string, error) {
	// Sanitize and clamp input
	sanitized := ai.SanitizeUserMessage(rawMessage, 400)
	if limit <= 0 || limit > 10 {
		limit = 5
	}

	// Cache lookup
	if out, ok := s.getFromCache(sanitized); ok {
		return out, "cached", nil
	}

	// Use web-style query for relevance
	products, _, err := s.repo.SearchProducts(ctx, sanitized, limit, 0, "relevance")
	if err != nil {
		return nil, "", err
	}

	// Build suggestions
	suggestions := make([]dto.RecommendSuggestion, 0, len(products))
	for _, p := range products {
		reason := buildSimpleReason(sanitized, []string(p.Occasions), []string(p.Seasons), []string(p.Accords))
		suggestions = append(suggestions, dto.RecommendSuggestion{
			ID:           p.ID,
			Name:         p.Name,
			Brand:        p.Brand,
			Slug:         p.Slug,
			ThumbnailURL: p.ThumbnailURL,
			Price:        p.Price,
			Reason:       reason,
		})
	}

	s.setCache(sanitized, suggestions)
	return suggestions, "fts", nil
}

func (s *AIService) getFromCache(key string) ([]dto.RecommendSuggestion, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ent, ok := s.cache[key]; ok {
		if time.Now().Before(ent.expiresAt) {
			return ent.suggestions, true
		}
		delete(s.cache, key)
	}
	return nil, false
}

func (s *AIService) setCache(key string, val []dto.RecommendSuggestion) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = cacheEntry{suggestions: val, expiresAt: time.Now().Add(s.ttl)}
}

// Very lightweight explanation by echoing matched keywords.
func buildSimpleReason(msg string, occasions, seasons, accords []string) string {
	// Lowercase tokens from message
	m := strings.ToLower(msg)
	b := &strings.Builder{}
	wrote := false
	if anyOverlap(m, occasions) {
		b.WriteString("Combina com sua ocasião mencionada")
		wrote = true
	}
	if anyOverlap(m, seasons) {
		if wrote {
			b.WriteString(" · ")
		}
		b.WriteString("Apropriado para a estação citada")
		wrote = true
	}
	if anyOverlap(m, accords) {
		if wrote {
			b.WriteString(" · ")
		}
		b.WriteString("Perfil olfativo alinhado (accords)")
		wrote = true
	}
	if !wrote {
		return "Relevância por similaridade de texto"
	}
	return b.String()
}

func anyOverlap(m string, arr []string) bool {
	for _, v := range arr {
		if v == "" {
			continue
		}
		if strings.Contains(m, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

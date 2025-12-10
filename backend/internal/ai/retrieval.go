package ai

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/embeddings"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// RetrievalService handles product retrieval with caching.
type RetrievalService struct {
	products repository.ProductRepository
	emb      embeddings.Provider

	mu    sync.Mutex
	cache map[string]cacheEntry
	ttl   time.Duration
}

type cacheEntry struct {
	suggestions []dto.RecommendSuggestion
	expiresAt   time.Time
}

// NewRetrievalService creates a new retrieval service.
func NewRetrievalService(repo repository.ProductRepository, embProvider embeddings.Provider) *RetrievalService {
	return &RetrievalService{
		products: repo,
		emb:      embProvider,
		cache:    make(map[string]cacheEntry),
		ttl:      5 * time.Minute,
	}
}

// GetSuggestions retrieves product suggestions using hybrid search.
func (r *RetrievalService) GetSuggestions(ctx context.Context, prefs Slots, msg string) []dto.RecommendSuggestion {
	msg = SanitizeUserMessage(msg, 0)
	q := strings.TrimSpace(BuildSearchQuery(prefs, msg))
	key := ProfileHash(prefs) + ":" + queryHash(q)
	if out, ok := r.getFromCache(key); ok {
		return out
	}

	// If embeddings provider is not available, fallback to pure FTS
	if r.emb == nil {
		sugs := []dto.RecommendSuggestion{}
		q := BuildSearchQuery(prefs, msg)
		q = strings.TrimSpace(q)
		if q != "" {
			gender := getGenderFilter(prefs)
			prods, _, _ := r.products.SearchProductsByGender(ctx, q, 3, 0, "relevance", gender)
			for _, p := range prods {
				reason := shortReason(prefs, p)
				sugs = append(sugs, dto.RecommendSuggestion{
					ID: p.ID, Name: p.Name, Brand: p.Brand, Slug: p.Slug, ThumbnailURL: p.ThumbnailURL, Price: p.Price,
					Reason: reason,
				})
			}
		}
		r.setCache(key, sugs)
		return sugs
	}

	// Parallel retrieval: FTS, Embeddings, and direct slot matching
	topK := 3
	acc := make([]dto.RecommendSuggestion, 0, topK*3)
	seen := make(map[uint]bool)

	type result struct {
		sugs []dto.RecommendSuggestion
		err  error
	}
	results := make(chan result, 3)

	// 1. FTS
	go func() {
		sugs := []dto.RecommendSuggestion{}
		q := BuildSearchQuery(prefs, msg)
		q = strings.TrimSpace(q)
		if q != "" {
			gender := getGenderFilter(prefs)
			prods, _, _ := r.products.SearchProductsByGender(ctx, q, topK, 0, "relevance", gender)
			for _, p := range prods {
				reason := shortReason(prefs, p)
				sugs = append(sugs, dto.RecommendSuggestion{
					ID: p.ID, Name: p.Name, Brand: p.Brand, Slug: p.Slug, ThumbnailURL: p.ThumbnailURL, Price: p.Price,
					Reason: reason,
				})
			}
		}
		results <- result{sugs: sugs}
	}()

	// 2. Embeddings
	go func() {
		sugs := []dto.RecommendSuggestion{}
		if r.emb != nil {
			queryText := BuildSearchQuery(prefs, msg)
			if queryText != "" {
				emb, err := r.emb.EmbedQuery(queryText)
				if err == nil && len(emb) > 0 {
					gender := getGenderFilter(prefs)
					similar, err := r.products.FindSimilarProductsByEmbeddingAndGender(ctx, emb, topK, gender)
					if err == nil {
						for _, p := range similar {
							reason := "Similaridade semântica com sua consulta"
							reason = shortReason(prefs, p) + " • " + reason
							sugs = append(sugs, dto.RecommendSuggestion{
								ID: p.ID, Name: p.Name, Brand: p.Brand, Slug: p.Slug, ThumbnailURL: p.ThumbnailURL, Price: p.Price,
								Reason: reason,
							})
						}
					}
				}
			}
		}
		results <- result{sugs: sugs}
	}()

	// 3. Direct slot matching
	go func() {
		sugs := []dto.RecommendSuggestion{}
		if len(prefs.Accords) > 0 {
			accordStr := strings.Join(prefs.Accords, " | ")
			q := fmt.Sprintf("(%s)", accordStr)
			gender := getGenderFilter(prefs)
			prods, _, _ := r.products.SearchProductsByGender(ctx, q, topK, 0, "relevance", gender)
			for _, p := range prods {
				reason := "Correspondência direta de acordes"
				reason = shortReason(prefs, p) + " • " + reason
				sugs = append(sugs, dto.RecommendSuggestion{
					ID: p.ID, Name: p.Name, Brand: p.Brand, Slug: p.Slug, ThumbnailURL: p.ThumbnailURL, Price: p.Price,
					Reason: reason,
				})
			}
		}
		results <- result{sugs: sugs}
	}()

	// Collect results
	for i := 0; i < 3; i++ {
		res := <-results
		if res.err == nil {
			for _, sug := range res.sugs {
				if !seen[sug.ID] {
					seen[sug.ID] = true
					acc = append(acc, sug)
				}
			}
		}
	}

	if len(acc) > topK {
		acc = acc[:topK]
	}

	r.setCache(key, acc)
	return acc
}

func queryHash(q string) string {
	if q == "" {
		q = "empty"
	}
	h := sha1.Sum([]byte(q))
	return hex.EncodeToString(h[:8])
}

func getGenderFilter(prefs Slots) string {
	if len(prefs.Gender) > 0 {
		return prefs.Gender[0]
	}
	return ""
}

func (r *RetrievalService) getFromCache(key string) ([]dto.RecommendSuggestion, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ent, ok := r.cache[key]; ok {
		if time.Now().Before(ent.expiresAt) {
			return ent.suggestions, true
		}
		delete(r.cache, key)
	}
	return nil, false
}

func (r *RetrievalService) setCache(key string, val []dto.RecommendSuggestion) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[key] = cacheEntry{suggestions: val, expiresAt: time.Now().Add(r.ttl)}
}

// ClearCache clears all cached suggestions.
func (r *RetrievalService) ClearCache() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache = make(map[string]cacheEntry)
}

func shortReason(p Slots, prod model.Product) string {
	lower := func(arr []string) []string {
		out := make([]string, 0, len(arr))
		for _, v := range arr {
			out = append(out, strings.ToLower(v))
		}
		return out
	}

	// Convert pq.StringArray to []string
	prodOccasions := make([]string, len(prod.Occasions))
	for i, v := range prod.Occasions {
		prodOccasions[i] = string(v)
	}
	prodSeasons := make([]string, len(prod.Seasons))
	for i, v := range prod.Seasons {
		prodSeasons[i] = string(v)
	}
	prodAccords := make([]string, len(prod.Accords))
	for i, v := range prod.Accords {
		prodAccords[i] = string(v)
	}

	score := 0
	score += overlapCount(lower(p.Occasions), lower(prodOccasions))
	score += overlapCount(lower(p.Seasons), lower(prodSeasons))
	score += overlapCount(lower(p.Accords), lower(prodAccords))
	if score == 0 {
		return "Compatível por perfil geral"
	}
	return fmt.Sprintf("%d ponto(s) de compatibilidade com suas preferências", score)
}

func overlapCount(a, b []string) int {
	m := map[string]bool{}
	for _, x := range a {
		m[x] = true
	}
	c := 0
	for _, y := range b {
		if m[y] {
			c++
		}
	}
	return c
}

package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/leoferamos/aroma-sense/internal/ai"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/embeddings"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/llm"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// ChatServiceInterface defines the minimal surface used by handlers for chat interactions.
type ChatServiceInterface interface {
	Chat(ctx context.Context, sessionID string, rawMsg string) (dto.ChatResponse, error)
	ClearRetrievalCache()
}

// ChatService orchestrates conversational recommendation with a lightweight LLM.
type ChatService struct {
	products  repository.ProductRepository
	llm       llm.Provider
	retrieval *ai.RetrievalService

	mu    sync.Mutex
	state map[string]*ai.Conversation
	ttl   time.Duration
}

func NewChatService(repo repository.ProductRepository, provider llm.Provider, embProvider embeddings.Provider) *ChatService {
	retrieval := ai.NewRetrievalService(repo, embProvider)
	return &ChatService{
		products:  repo,
		llm:       provider,
		retrieval: retrieval,
		state:     make(map[string]*ai.Conversation),
		ttl:       30 * time.Minute,
	}
}

// Chat processes a user message and returns an LLM reply plus optional product suggestions.
func (s *ChatService) Chat(ctx context.Context, sessionID string, rawMsg string) (dto.ChatResponse, error) {
	sid := sessionID
	if sid == "" {
		sid = "anon"
	}
	sanitized := ai.SanitizeUserMessage(rawMsg, 500)

	if isGreetingOnly(sanitized) {
		reply := "Olá! Eu sou a assistente da Aroma Sense. Posso te ajudar a escolher perfumes por ocasião, estação, acordes (cítrico, floral, amadeirado), intensidade e orçamento. Como você quer começar?"
		return dto.ChatResponse{Reply: reply, Suggestions: nil, FollowUpHint: "Prefere algo cítrico, floral ou amadeirado?"}, nil
	}
	if isFarewell(sanitized) {
		reply := "Até logo! Quando quiser, volto a te ajudar com recomendações de perfumes."
		return dto.ChatResponse{Reply: reply}, nil
	}

	conv := s.getOrCreate(sid)
	// Parse latest user message
	parsed := ai.Parse(sanitized)

	// Detect context reset: if user message contains a new main slot (occasion, season, climate, intensity) or reset keywords, overwrite those slots and clear others
	resetKeywords := []string{"agora", "na verdade", "quero outro", "outro", "mudei de ideia", "dessa vez", "desta vez", "novo pedido", "diferente", "mudou", "trocar", "quero"}
	msgLower := strings.ToLower(sanitized)
	reset := false
	for _, kw := range resetKeywords {
		if strings.Contains(msgLower, kw) {
			reset = true
			break
		}
	}
	// If user provided a new occasion, season, climate, or intensity, treat as context reset
	if len(parsed.Occasions) > 0 || len(parsed.Seasons) > 0 || len(parsed.Climate) > 0 || len(parsed.Intensity) > 0 {
		reset = true
	}
	if reset {
		// Overwrite main slots, clear others
		conv.Prefs.Occasions = parsed.Occasions
		conv.Prefs.Seasons = parsed.Seasons
		conv.Prefs.Climate = parsed.Climate
		conv.Prefs.Intensity = parsed.Intensity
		conv.Prefs.Accords = parsed.Accords
		conv.Prefs.Budget = parsed.Budget
		conv.Prefs.Longevity = parsed.Longevity
		conv.Prefs.Gender = parsed.Gender
		conv.Prefs.Notes = parsed.Notes
	} else {
		conv.Prefs = ai.Merge(conv.Prefs, parsed)
	}
	conv.AddMessage(sanitized, parsed)

	// If message seems off-topic and we don't have preferences yet
	if !isOnTopic(sanitized) && isEmptyPrefs(conv.Prefs) {
		reply := "Posso ajudar especificamente com perfumes. Me diga, por exemplo: ocasião (trabalho, festa), acordes que curte (cítrico, floral, amadeirado), intensidade (suave, moderada, forte) e sua faixa de preço."
		follow := ai.BuildFollowUpHint(conv.Prefs)
		return dto.ChatResponse{Reply: reply, FollowUpHint: follow}, nil
	}

	// Only suggest products if user has provided enough preferences
	if !hasMinimumPrefs(conv.Prefs) {
		reply := "Para te sugerir perfumes, preciso saber um pouco mais sobre suas preferências. Me diga, por exemplo: ocasião (trabalho, festa), acordes que curte (cítrico, floral, amadeirado), intensidade (suave, moderada, forte) ou gênero."
		follow := ai.BuildFollowUpHint(conv.Prefs)
		return dto.ChatResponse{Reply: reply, Suggestions: nil, FollowUpHint: follow}, nil
	}

	// Retrieve candidate products
	suggestions := s.retrieval.GetSuggestions(ctx, conv.Prefs, sanitized)

	prompt := ai.BuildPrompt(conv, sanitized, suggestions)
	reply, err := s.llm.Generate(ctx, prompt, 180)
	if err != nil {
		// Fallback deterministic reply while still returning suggestions
		names := make([]string, 0, len(suggestions))
		for i := 0; i < len(suggestions) && i < 3; i++ {
			names = append(names, suggestions[i].Name)
		}
		base := "Tenho algumas sugestões"
		if len(names) > 0 {
			base += ": " + strings.Join(names, ", ")
		}
		reply = base
	}

	// Return up to 2 suggestions
	var limitedSuggestions []dto.RecommendSuggestion
	if len(suggestions) > 0 {
		if len(suggestions) > 2 {
			limitedSuggestions = suggestions[:2]
		} else {
			limitedSuggestions = suggestions
		}
		mentioned := false
		for _, s := range suggestions {
			if strings.Contains(strings.ToLower(reply), strings.ToLower(s.Name)) {
				mentioned = true
				break
			}
		}
		if !mentioned {
			reply += "\n\nBaseado no que você disse, tenho uma sugestão: " + suggestions[0].Name + " da " + suggestions[0].Brand + "."
			if len(limitedSuggestions) > 1 {
				reply += " Ou " + limitedSuggestions[1].Name + " da " + limitedSuggestions[1].Brand + "."
			}
		}
		conv.LastSuggestions = true
	}

	var follow string
	if len(suggestions) > 0 || conv.TurnCount > 2 || hasMinimumPrefs(conv.Prefs) {
		follow = ""
	} else {
		follow = ai.BuildFollowUpHint(conv.Prefs)
	}
	return dto.ChatResponse{Reply: reply, Suggestions: limitedSuggestions, FollowUpHint: follow}, nil
}

func (s *ChatService) getOrCreate(id string) *ai.Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.state[id]; ok {
		if time.Now().Before(c.ExpiresAt) {
			return c
		}
	}
	c := ai.NewConversation()
	s.state[id] = c
	return c
}

// --- lightweight intent/topic helpers ---
func isGreetingOnly(s string) bool {
	t := strings.TrimSpace(strings.ToLower(s))
	if t == "" {
		return false
	}
	// common greetings in pt/en with optional punctuation
	greetings := []string{"oi", "ola", "olá", "hello", "hi", "hey", "bom dia", "boa tarde", "boa noite"}
	for _, g := range greetings {
		if t == g || t == g+"!" || t == g+"." {
			return true
		}
	}
	return false
}

func isFarewell(s string) bool {
	t := strings.TrimSpace(strings.ToLower(s))
	farewells := []string{"tchau", "até", "até mais", "valeu", "obrigado", "obrigada"}
	for _, f := range farewells {
		if strings.HasPrefix(t, f) {
			return true
		}
	}
	return false
}

func isOnTopic(s string) bool {
	t := strings.ToLower(s)
	// simple keyword gate for fragrance domain
	keys := []string{"perfume", "fragr", "cheiro", "aroma", "odor", "eau", "parfum", "toilette"}
	for _, k := range keys {
		if strings.Contains(t, k) {
			return true
		}
	}
	return false
}

func (s *ChatService) ClearRetrievalCache() {
	s.retrieval.ClearCache()
}

func isEmptyPrefs(p ai.Slots) bool {
	return len(p.Occasions) == 0 && len(p.Seasons) == 0 && len(p.Accords) == 0 && len(p.Intensity) == 0 && len(p.Climate) == 0 && len(p.Budget) == 0 && len(p.Longevity) == 0
}

func hasMinimumPrefs(p ai.Slots) bool {
	count := 0
	if len(p.Occasions) > 0 {
		count++
	}
	if len(p.Accords) > 0 {
		count++
	}
	if len(p.Gender) > 0 {
		count++
	}
	if len(p.Intensity) > 0 {
		count++
	}
	return count >= 2
}

package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/leoferamos/aroma-sense/internal/ai"
	"github.com/leoferamos/aroma-sense/internal/ai/slots"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/embeddings"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/llm"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// ChatService orchestrates conversational recommendation with a lightweight LLM.
type ChatService struct {
	products  repository.ProductRepository
	llm       llm.Provider
	retrieval *ai.RetrievalService

	mu    sync.Mutex
	state map[string]*ai.Conversation
	ttl   time.Duration
}

func NewChatService(repo repository.ProductRepository, provider llm.Provider, embProvider embeddings.Provider, aiService *AIService) *ChatService {
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

	// Lightweight NLU shortcuts
	if isGreetingOnly(sanitized) {
		reply := "Olá! Eu sou a assistente da Aroma Sense. Posso te ajudar a escolher perfumes por ocasião, estação, acordes (cítrico, floral, amadeirado), intensidade e orçamento. Como você quer começar?"
		return dto.ChatResponse{Reply: reply, Suggestions: nil, FollowUpHint: "Prefere algo cítrico, floral ou amadeirado?"}, nil
	}
	if isFarewell(sanitized) {
		reply := "Até logo! Quando quiser, volto a te ajudar com recomendações de perfumes."
		return dto.ChatResponse{Reply: reply}, nil
	}

	conv := s.getOrCreate(sid)
	// Update conversation state
	conv.AddMessage(sanitized, slots.Parse(sanitized))

	// If message seems off-topic and we don't have preferences yet
	if !isOnTopic(sanitized) && isEmptyPrefs(conv.Prefs) {
		reply := "Posso ajudar especificamente com perfumes. Me diga, por exemplo: ocasião (trabalho, festa), acordes que curte (cítrico, floral, amadeirado), intensidade (suave, moderada, forte) e sua faixa de preço."
		follow := ai.BuildFollowUpHint(conv.Prefs)
		return dto.ChatResponse{Reply: reply, FollowUpHint: follow}, nil
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

	// Append suggestions if not mentioned in reply
	if len(suggestions) > 0 {
		mentioned := false
		for _, s := range suggestions {
			if strings.Contains(strings.ToLower(reply), strings.ToLower(s.Name)) {
				mentioned = true
				break
			}
		}
		if !mentioned {
			reply += "\n\nBaseado no que você disse, tenho uma sugestão: " + suggestions[0].Name + " da " + suggestions[0].Brand + "."
			if len(suggestions) > 1 {
				reply += " Ou " + suggestions[1].Name + " da " + suggestions[1].Brand + "."
			}
		}
	}

	follow := ai.BuildFollowUpHint(conv.Prefs)
	return dto.ChatResponse{Reply: reply, Suggestions: suggestions, FollowUpHint: follow}, nil
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

func isEmptyPrefs(p slots.Slots) bool {
	return len(p.Occasions) == 0 && len(p.Seasons) == 0 && len(p.Accords) == 0 && len(p.Intensity) == 0 && len(p.Climate) == 0 && len(p.Budget) == 0 && len(p.Longevity) == 0
}

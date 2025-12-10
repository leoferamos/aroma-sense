package ai

import (
	"strings"
	"time"
)

// Conversation represents a chat session state.
type Conversation struct {
	History         []string
	Summary         string
	Prefs           Slots
	ExpiresAt       time.Time
	TurnCount       int
	LastSuggestions bool
}

// NewConversation creates a new conversation.
func NewConversation() *Conversation {
	return &Conversation{
		History:   []string{},
		Prefs:     Slots{},
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
}

// AddMessage adds a user message and updates preferences.
func (c *Conversation) AddMessage(msg string, prefs Slots) {
	msg = SanitizeUserMessage(msg, 0)
	c.Prefs = Merge(c.Prefs, prefs)
	c.History = append(c.History, msg)
	c.TurnCount++
	if len(c.History) > 12 {
		c.summarize()
	}
	c.ExpiresAt = time.Now().Add(30 * time.Minute)
}

// summarize creates a summary when history is too long.
func (c *Conversation) summarize() {
	if len(c.History) < 6 {
		return
	}
	prefParts := []string{}
	if len(c.Prefs.Occasions) > 0 {
		prefParts = append(prefParts, "Ocasiões="+strings.Join(c.Prefs.Occasions, ","))
	}
	if len(c.Prefs.Climate) > 0 {
		prefParts = append(prefParts, "Clima="+strings.Join(c.Prefs.Climate, ","))
	}
	if len(c.Prefs.Intensity) > 0 {
		prefParts = append(prefParts, "Intensidade="+strings.Join(c.Prefs.Intensity, ","))
	}
	if len(c.Prefs.Accords) > 0 {
		prefParts = append(prefParts, "Acordes="+strings.Join(c.Prefs.Accords, ","))
	}
	if len(c.Prefs.Budget) > 0 {
		prefParts = append(prefParts, "Orçamento="+strings.Join(c.Prefs.Budget, ","))
	}
	if len(c.Prefs.Longevity) > 0 {
		prefParts = append(prefParts, "Longevidade="+strings.Join(c.Prefs.Longevity, ","))
	}
	if len(c.Prefs.Seasons) > 0 {
		prefParts = append(prefParts, "Estações="+strings.Join(c.Prefs.Seasons, ","))
	}
	if len(c.Prefs.Gender) > 0 {
		prefParts = append(prefParts, "Gênero="+strings.Join(c.Prefs.Gender, ","))
	}
	if len(c.Prefs.Notes) > 0 {
		prefParts = append(prefParts, "Notas="+strings.Join(c.Prefs.Notes, ","))
	}
	c.Summary = strings.Join(prefParts, " | ")
	last := c.History[len(c.History)-2:]
	c.History = append([]string{"(resumo) " + c.Summary}, last...)
}

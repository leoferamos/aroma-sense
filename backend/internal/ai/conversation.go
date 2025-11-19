package ai

import (
	"strings"
	"time"

	"github.com/leoferamos/aroma-sense/internal/ai/slots"
)

// Conversation represents a chat session state.
type Conversation struct {
	History   []string
	Summary   string
	Prefs     slots.Slots
	ExpiresAt time.Time
}

// NewConversation creates a new conversation.
func NewConversation() *Conversation {
	return &Conversation{
		History:   []string{},
		Prefs:     slots.Slots{},
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
}

// AddMessage adds a user message and updates preferences.
func (c *Conversation) AddMessage(msg string, prefs slots.Slots) {
	c.Prefs = slots.Merge(c.Prefs, prefs)
	c.History = append(c.History, msg)
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
	if len(c.Prefs.Seasons) > 0 {
		prefParts = append(prefParts, "Estações="+strings.Join(c.Prefs.Seasons, ","))
	}
	if len(c.Prefs.Accords) > 0 {
		prefParts = append(prefParts, "Accords="+strings.Join(c.Prefs.Accords, ","))
	}
	c.Summary = strings.Join(prefParts, " | ")
	last := c.History[len(c.History)-2:]
	c.History = append([]string{"(resumo) " + c.Summary}, last...)
}

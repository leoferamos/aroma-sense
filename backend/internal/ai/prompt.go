package ai

import (
	"strings"

	"github.com/leoferamos/aroma-sense/internal/dto"
)

// BuildPrompt constructs the LLM prompt for chat response.
func BuildPrompt(c *Conversation, userMsg string, suggestions []dto.RecommendSuggestion) string {
	var sb strings.Builder
	sb.WriteString("Você é um assistente de perfumes amigável. Responda em português de forma curta e natural. Colete preferências perguntando educadamente se faltar info.\n")
	if c.Summary != "" {
		sb.WriteString("Preferências: " + c.Summary + "\n")
	}
	sb.WriteString("Mensagem do usuário: " + userMsg + "\n")
	return sb.String()
}

// BuildFollowUpHint returns a hint for missing slots.
func BuildFollowUpHint(p Slots) string {
	missing := NextMissing(p)
	if missing == "Accords" {
		return "Você tem preferência por algum acorde (cítrico, floral, amadeirado)?"
	}
	if missing == "Occasions" {
		return "Vai usar mais em qual ocasião (trabalho, festa, encontro)?"
	}
	if missing == "Seasons" {
		return "Alguma estação específica (verão, inverno)?"
	}
	if missing == "Intensity" {
		return "Prefere algo suave, moderado ou forte?"
	}
	if missing == "Climate" {
		return "O clima será mais frio, quente, úmido ou seco?"
	}
	if missing == "Budget" {
		return "Tem alguma faixa de preço em mente?"
	}
	return "Posso refinar por intensidade, preço ou longevidade."
}

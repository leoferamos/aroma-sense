package ai

import (
	"strings"

	"github.com/leoferamos/aroma-sense/internal/ai/slots"
	"github.com/leoferamos/aroma-sense/internal/dto"
)

// BuildPrompt constructs the LLM prompt for chat response.
func BuildPrompt(c *Conversation, userMsg string, suggestions []dto.RecommendSuggestion) string {
	var sb strings.Builder
	sb.WriteString("Você é um assistente de perfumes amigável. Responda APENAS com uma mensagem conversacional curta em português. Não repita a pergunta do usuário. Não liste produtos - apenas converse naturalmente e mencione sugestões se apropriado.\n")
	if c.Summary != "" {
		sb.WriteString("Preferências do usuário: " + c.Summary + "\n")
	}
	sb.WriteString("Mensagem do usuário: " + userMsg + "\n")
	if len(suggestions) > 0 {
		sb.WriteString("Sugestões disponíveis (use apenas se relevante para a conversa):\n")
		max := len(suggestions)
		if max > 3 {
			max = 3
		}
		for i := 0; i < max; i++ {
			s := suggestions[i]
			sb.WriteString("- " + s.Name + " (" + s.Brand + ") - " + s.Reason + "\n")
		}
	}
	missing := slots.NextMissing(c.Prefs)
	if missing != "" {
		hint := BuildFollowUpHint(c.Prefs)
		sb.WriteString("Como não temos info sobre " + strings.ToLower(missing) + ", pergunte educadamente: " + hint + "\n")
	}
	sb.WriteString("Responda de forma natural, curta e útil.\n")
	return sb.String()
}

// BuildFollowUpHint returns a hint for missing slots.
func BuildFollowUpHint(p slots.Slots) string {
	if slots.NextMissing(p) == "Accords" {
		return "Você tem preferência por algum acorde (cítrico, floral, amadeirado)?"
	}
	if slots.NextMissing(p) == "Occasions" {
		return "Vai usar mais em qual ocasião (trabalho, festa, encontro)?"
	}
	if slots.NextMissing(p) == "Seasons" {
		return "Alguma estação específica (verão, inverno)?"
	}
	if slots.NextMissing(p) == "Intensity" {
		return "Prefere algo suave, moderado ou forte?"
	}
	if slots.NextMissing(p) == "Climate" {
		return "O clima será mais frio, quente, úmido ou seco?"
	}
	if slots.NextMissing(p) == "Budget" {
		return "Tem alguma faixa de preço em mente?"
	}
	return "Posso refinar por intensidade, preço ou longevidade."
}

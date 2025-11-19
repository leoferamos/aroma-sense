package ai

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"
)

// Slots represents the user's preference state captured over the conversation.
type Slots struct {
	Occasions []string
	Climate   []string
	Seasons   []string
	Intensity []string
	Accords   []string
	Budget    []string
	Longevity []string
}

// Order of slot clarification
var order = []string{"Occasions", "Climate", "Intensity", "Accords", "Budget", "Longevity", "Seasons"}

// Keywords for simple extraction
var (
	occasionKw  = []string{"noite", "dia", "trabalho", "casual", "encontro", "festa"}
	seasonKw    = []string{"verão", "inverno", "outono", "primavera"}
	climateKw   = []string{"frio", "quente", "úmido", "seco"}
	intensityKw = []string{"suave", "leve", "moderad", "forte", "intens"}
	accordKw    = []string{"cítrico", "citric", "amadeir", "floral", "oriental", "especiaria", "especiad", "frutado", "verde", "baunilha", "musk"}
	budgetKw    = []string{"barato", "baixo", "médio", "medio", "alto", "caro"}
	longevityKw = []string{"curta", "média", "media", "longa", "fixação", "durabilidade"}
)

// Parse extracts slot values from a single message.
func Parse(msg string) Slots {
	m := strings.ToLower(msg)
	s := Slots{}
	s.Occasions = matchAny(m, occasionKw)
	s.Seasons = matchAny(m, seasonKw)
	s.Climate = matchAny(m, climateKw)
	s.Intensity = matchAny(m, intensityKw)
	s.Accords = matchAny(m, accordKw)
	s.Budget = matchAny(m, budgetKw)
	s.Longevity = matchAny(m, longevityKw)
	return s
}

// Merge returns a new Slots combining existing with new values (dedup, stable order).
func Merge(a, b Slots) Slots {
	out := Slots{
		Occasions: dedup(append([]string{}, append(a.Occasions, b.Occasions...)...)),
		Climate:   dedup(append([]string{}, append(a.Climate, b.Climate...)...)),
		Seasons:   dedup(append([]string{}, append(a.Seasons, b.Seasons...)...)),
		Intensity: dedup(append([]string{}, append(a.Intensity, b.Intensity...)...)),
		Accords:   dedup(append([]string{}, append(a.Accords, b.Accords...)...)),
		Budget:    dedup(append([]string{}, append(a.Budget, b.Budget...)...)),
		Longevity: dedup(append([]string{}, append(a.Longevity, b.Longevity...)...)),
	}
	return out
}

// NextMissing returns the next slot name to clarify according to the defined order.
func NextMissing(s Slots) string {
	for _, name := range order {
		switch name {
		case "Occasions":
			if len(s.Occasions) == 0 {
				return name
			}
		case "Climate":
			if len(s.Climate) == 0 {
				return name
			}
		case "Intensity":
			if len(s.Intensity) == 0 {
				return name
			}
		case "Accords":
			if len(s.Accords) == 0 {
				return name
			}
		case "Budget":
			if len(s.Budget) == 0 {
				return name
			}
		case "Longevity":
			if len(s.Longevity) == 0 {
				return name
			}
		case "Seasons":
			if len(s.Seasons) == 0 {
				return name
			}
		}
	}
	return ""
}

// ProfileHash returns a deterministic short hash of the filled slots for caching.
func ProfileHash(s Slots) string {
	parts := []string{
		join("oc", s.Occasions),
		join("cl", s.Climate),
		join("in", s.Intensity),
		join("ac", s.Accords),
		join("bu", s.Budget),
		join("lo", s.Longevity),
		join("se", s.Seasons),
	}
	h := sha1.Sum([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(h[:8])
}

// BuildSearchQuery combines user message with top preferences to full-text search.
func BuildSearchQuery(s Slots, msg string) string {
	parts := []string{msg}
	if len(s.Accords) > 0 {
		parts = append(parts, s.Accords[0])
	}
	if len(s.Occasions) > 0 {
		parts = append(parts, s.Occasions[0])
	}
	if len(s.Seasons) > 0 {
		parts = append(parts, s.Seasons[0])
	}
	if len(s.Climate) > 0 {
		parts = append(parts, s.Climate[0])
	}
	return strings.Join(parts, " ")
}

func matchAny(m string, kws []string) []string {
	out := make([]string, 0)
	for _, k := range kws {
		if strings.Contains(m, k) {
			out = append(out, k)
		}
	}
	return dedup(out)
}

func dedup(arr []string) []string {
	if len(arr) == 0 {
		return arr
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	out := arr[:0]
	var prev string
	for i, v := range arr {
		if i == 0 || v != prev {
			out = append(out, v)
			prev = v
		}
	}
	return out
}

func join(prefix string, arr []string) string {
	if len(arr) == 0 {
		return prefix + ":"
	}
	return prefix + ":" + strings.Join(arr, ",")
}

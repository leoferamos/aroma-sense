package ai

import (
	"crypto/sha1"
	"encoding/hex"
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
	Gender    []string
	Notes     []string
}

// Order of slot clarification
var order = []string{"Occasions", "Climate", "Intensity", "Accords", "Budget", "Longevity", "Seasons", "Gender"}

// Keywords for simple extraction
var (
	occasionKw = []string{
		"dia a dia", "cotidiano", "dia",
		"trabalho", "escritório", "office",
		"academia", "treino", "corrida", "exercício", "exercicio", "esporte",
		"pós-banho", "pos banho", "banho",
		"festa", "balada", "rolê", "role",
		"encontro", "date",
		"casual",
		"noite", "noturno",
		"formal", "evento", "evento especial", "eventos especiais",
		// expressões do seu catálogo
		"noite casual", "casual noturno", "casual chique",
	}

	seasonKw = []string{
		"verão", "inverno", "outono", "primavera",
		// expressões do seu catálogo
		"noites de inverno", "noites de primavera",
	}

	climateKw = []string{
		"frio", "quente", "úmido", "seco",
	}

	intensityKw = []string{
		"discreto",
		"suave", "leve",
		"moderado", "moderad",
		"forte",
		"intens", "intenso", "muito intenso",
		"marcante",
		"projeção", "projecao",
		"rastro",
	}

	accordKw = []string{
		// principais famílias
		"aromático", "aromatic",
		"aquático", "aquatico", "aquatic",
		"marinho", "nota marinha", "notas marinhas",
		"fresco", "fresh",
		"limpo", "clean",

		"âmbar", "amber", "ambar",
		"incenso", "incensado",
		"almiscarado", "almíscar", "almiscar", "musk",
		"amadeir", "woody",
		"floral",
		"frutado",
		"verde",
		"cítrico", "citrus", "citric",

		"baunilha", "vanilla",
		"doce", "adocicado",
		"atalcado",
		"oriental",

		// especiados do catálogo
		"especiado fresco", "especiado quente",
		"especiad", "especiaria", "especiado",

		// acordes/descritores que aparecem nos seus produtos
		"lavanda",
		"terroso",
		"esfumaçado", "esfumacado",
		"floral branco",
		"chá", "cha",
		"espumante",
		"café", "cafe",
		"íris", "iris",

		// termo comum ligado ao seu portfólio (ex.: Black Opium)
		"gourmand",
	}

	budgetKw = []string{
		"barato", "baixo",
		"acessível", "acessivel",
		"médio", "medio",
		"custo-benefício", "custo beneficio",
		"alto", "caro",
		"premium", "luxo",
	}

	longevityKw = []string{
		"curta",
		"média", "media",
		"longa",
		"fixação", "fixacao",
		"durabilidade",
		"fixa bem",
		"dura bastante",
		"performance",
	}

	genderKw = []string{
		"masculino", "masc", "homem",
		"feminino", "fem", "mulher",
		"unissex",
	}

	notesKw = []string{
		"lavanda",
		"bergamota",
		"limão", "limao",
		"mandarina",
		"hortelã", "hortela",
		"cardamomo",
		"patchouli", "patchuli",
		"incenso",
		"âmbar", "ambar",
		"almíscar", "almiscar", "musk",
		"ambroxan",
		"vetiver",
		"jasmim",
		"flor de laranjeira",
		"íris", "iris",
		"canela",
		"chá verde", "cha verde", "chá", "cha",
		"pera",
		"maçã", "maca",
		"toranja",
		"notas marinhas", "nota marinha",
	}
)

// Parse extracts slot values from a single message.
func Parse(msg string) Slots {
	m := strings.ToLower(msg)
	s := Slots{}
	s.Occasions = matchAny(m, occasionKw, "occasion")
	s.Seasons = matchAny(m, seasonKw, "season")
	s.Climate = matchAny(m, climateKw, "climate")
	s.Intensity = matchAny(m, intensityKw, "intensity")
	s.Accords = matchAny(m, accordKw, "accord")
	s.Budget = matchAny(m, budgetKw, "budget")
	s.Longevity = matchAny(m, longevityKw, "longevity")
	s.Gender = matchAny(m, genderKw, "gender")
	s.Notes = matchAny(m, notesKw, "note")
	return s
}

// Merge returns a new Slots combining existing with new values (dedup, stable order).
func Merge(a, b Slots) Slots {
	out := Slots{
		Occasions: dedupPreserveOrder(append([]string{}, append(a.Occasions, b.Occasions...)...)),
		Climate:   dedupPreserveOrder(append([]string{}, append(a.Climate, b.Climate...)...)),
		Seasons:   dedupPreserveOrder(append([]string{}, append(a.Seasons, b.Seasons...)...)),
		Intensity: dedupPreserveOrder(append([]string{}, append(a.Intensity, b.Intensity...)...)),
		Accords:   dedupPreserveOrder(append([]string{}, append(a.Accords, b.Accords...)...)),
		Budget:    dedupPreserveOrder(append([]string{}, append(a.Budget, b.Budget...)...)),
		Longevity: dedupPreserveOrder(append([]string{}, append(a.Longevity, b.Longevity...)...)),
		Gender:    dedupPreserveOrder(append([]string{}, append(a.Gender, b.Gender...)...)),
		Notes:     dedupPreserveOrder(append([]string{}, append(a.Notes, b.Notes...)...)),
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
		case "Gender":
			if len(s.Gender) == 0 {
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
		join("gn", s.Gender),
		join("nt", s.Notes),
	}
	h := sha1.Sum([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(h[:8])
}

// BuildSearchQuery combines user message with top preferences to full-text search.
func BuildSearchQuery(s Slots, msg string) string {
	parts := []string{msg}
	if len(s.Gender) > 0 {
		parts = append(parts, s.Gender[0])
	}
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
	if len(s.Notes) > 0 {
		parts = append(parts, s.Notes[0])
	}
	return strings.Join(parts, " ")
}

func matchAny(m string, kws []string, slotType string) []string {
	out := make([]string, 0)
	for _, k := range kws {
		if strings.Contains(m, k) {
			normalized := normalizeSlotValue(k, slotType)
			out = append(out, normalized)
		}
	}
	return dedupPreserveOrder(out)
}

func dedupPreserveOrder(arr []string) []string {
	if len(arr) == 0 {
		return arr
	}
	seen := make(map[string]bool)
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func normalizeSlotValue(v string, slotType string) string {
	v = strings.ToLower(strings.TrimSpace(v))

	if slotType == "gender" {
		switch v {
		case "masculino", "masc", "homem":
			return "Masculino"
		case "feminino", "fem", "mulher":
			return "Feminino"
		case "unissex":
			return "Unissex"
		}
	}

	if slotType == "accord" {
		switch v {
		case "aromático", "aromatic":
			return "Aromático"
		case "aquático", "aquatico", "aquatic", "marinho", "nota marinha":
			return "Aquático"
		case "fresco", "fresh", "limpo", "clean":
			return "Fresco"
		case "âmbar", "amber", "ambar":
			return "Âmbar"
		case "incenso", "incensado":
			return "Incensado"
		case "almiscarado", "almíscar", "almiscar", "musk":
			return "Almiscarado"
		case "amadeir", "woody":
			return "Amadeirado"
		case "doce", "adocicado":
			return "Doce"
		case "atalcado":
			return "Atalcado"
		case "cítrico", "citrus", "citric":
			return "Cítrico"
		case "floral":
			return "Floral"
		case "frutado":
			return "Frutado"
		case "verde":
			return "Verde"
		case "baunilha":
			return "Baunilha"
		case "especiado fresco":
			return "Especiado Fresco"
		case "especiado quente":
			return "Especiado Quente"
		case "especiad", "especiaria", "especiado":
			return "Especiado"
		case "oriental":
			return "Oriental"
		}
	}

	if slotType == "occasion" {
		switch v {
		case "dia a dia", "cotidiano":
			return "Dia a Dia"
		case "trabalho", "escritório", "office":
			return "Trabalho"
		case "academia", "treino", "corrida", "exercício", "exercicio", "esporte":
			return "Academia"
		case "festa", "balada", "rolê", "role":
			return "Festa"
		case "encontro", "date":
			return "Encontro"
		case "pós-banho", "pos banho", "banho":
			return "Pós-banho"
		case "noite", "noturno":
			return "Noturno"
		case "formal", "evento", "evento especial":
			return "Formal"
		case "casual":
			return "Casual"
		}
	}

	if slotType == "intensity" {
		switch v {
		case "discreto":
			return "Discreto"
		case "suave", "leve":
			return "Suave"
		case "moderado", "moderad":
			return "Moderado"
		case "forte", "intens", "intenso", "marcante", "projeção", "projecao":
			return "Forte"
		case "rastro":
			return "Rastro"
		}
	}

	if slotType == "budget" {
		switch v {
		case "barato", "baixo", "acessível", "acessivel":
			return "Acessível"
		case "médio", "medio", "custo-benefício", "custo beneficio":
			return "Médio"
		case "alto", "caro", "premium", "luxo":
			return "Premium"
		}
	}

	if slotType == "longevity" {
		switch v {
		case "curta":
			return "Curta"
		case "média", "media":
			return "Média"
		case "longa", "fixa bem", "dura bastante", "performance":
			return "Longa"
		case "fixação", "durabilidade":
			return "Durabilidade"
		}
	}

	if slotType == "note" {
		switch v {
		case "lavanda":
			return "Lavanda"
		case "bergamota":
			return "Bergamota"
		case "limão", "limao":
			return "Limão"
		case "mandarina":
			return "Mandarina"
		case "hortelã", "hortela":
			return "Hortelã"
		case "cardamomo":
			return "Cardamomo"
		case "patchouli", "patchuli":
			return "Patchouli"
		case "incenso":
			return "Incenso"
		case "âmbar", "ambar":
			return "Âmbar"
		case "almíscar", "almiscar", "musk":
			return "Almíscar"
		case "ambroxan":
			return "Ambroxan"
		case "vetiver":
			return "Vetiver"
		case "jasmim":
			return "Jasmim"
		case "flor de laranjeira":
			return "Flor de Laranjeira"
		case "íris", "iris":
			return "Íris"
		case "canela":
			return "Canela"
		case "chá verde", "cha verde", "chá", "cha":
			return "Chá Verde"
		case "pera":
			return "Pera"
		case "maçã", "maca":
			return "Maçã"
		case "toranja":
			return "Toranja"
		case "notas marinhas", "nota marinha":
			return "Notas Marinhas"
		}
	}

	return v
}

func dedup(arr []string) []string {
	if len(arr) == 0 {
		return arr
	}
	seen := make(map[string]bool)
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
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

package validation

import "regexp"

var cepDigits = regexp.MustCompile(`[^0-9]`)

// NormalizeCEP removes all non-digit characters from a CEP string.
func NormalizeCEP(s string) string {
	return cepDigits.ReplaceAllString(s, "")
}

// ExtractCEPFromString tries to extract a CEP from a free-form string.
func ExtractCEPFromString(s string) string {
	cleaned := NormalizeCEP(s)
	if len(cleaned) >= 8 {
		return cleaned[len(cleaned)-8:]
	}
	if len(cleaned) >= 5 {
		return cleaned
	}
	return ""
}

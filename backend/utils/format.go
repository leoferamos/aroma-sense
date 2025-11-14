package utils

import "strconv"

// FormatFloatTrim formats a float with fixed precision and trims trailing zeros and dot.
func FormatFloatTrim(v float64, prec int) string {
	s := strconv.FormatFloat(v, 'f', prec, 64)
	// trim trailing zeros and dot
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	if s == "" {
		return "0"
	}
	return s
}

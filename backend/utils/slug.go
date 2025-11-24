package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Slugify creates a URL-friendly slug from the provided parts.
func Slugify(parts ...string) string {
	base := strings.Join(parts, "-")
	base = strings.ToLower(strings.TrimSpace(base))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	noAccents, _, _ := transform.String(t, base)

	// Replace non-alphanumeric with hyphen
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := re.ReplaceAllString(noAccents, "-")

	// Collapse multiple hyphens
	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")

	// Trim hyphens
	slug = strings.Trim(slug, "-")

	// Guard empty result
	if slug == "" {
		slug = "product"
	}
	// Limit to 128 chars
	if len(slug) > 128 {
		slug = slug[:128]
		slug = strings.Trim(slug, "-")
	}
	return slug
}

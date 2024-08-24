package utils

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func ToSnakeCase(input string) string {
	// Remover acentos
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, _ := transform.String(t, input)

	var snake strings.Builder
	previousWasLower := false

	for i, r := range normalized {
		if unicode.IsUpper(r) {
			if i > 0 && previousWasLower {
				snake.WriteRune('_')
			}
			previousWasLower = false
		} else {
			previousWasLower = true
		}
		snake.WriteRune(unicode.ToLower(r))
	}

	reg := regexp.MustCompile(`[^a-z0-9]+`)
	snakeCase := reg.ReplaceAllString(snake.String(), "_")
	return strings.Trim(snakeCase, "_")
}

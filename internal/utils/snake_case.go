package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ToSnakeCase converte uma string para snake_case
func ToSnakeCase(input string) string {
	var snake strings.Builder
	previousWasLower := false

	for i, r := range input {
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

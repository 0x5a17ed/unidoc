package unidoc

import (
	"strings"
)

// translateMap translates a string using a rune mapping.
func translateMap(s string, m map[rune]rune) string {
	return strings.Map(func(inp rune) rune {
		if out, exists := m[inp]; exists {
			inp = out // Translate characters based on the mapping
		}
		return inp
	}, s)
}

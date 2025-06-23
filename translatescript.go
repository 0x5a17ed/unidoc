package unidoc

// toItalicScriptText converts regular text to mathematical script Unicode (𝒯𝒽𝒾𝓈)
func toItalicScriptText(text string) string {
	// Mathematical Script Unicode mapping
	m := map[rune]rune{
		'A': '\U0001D49C', 'B': '\U0001D4B7', 'C': '\U0001D49E', 'D': '\U0001D49F', 'E': '\U00002130',
		'F': '\U00002131', 'G': '\U0001D4A2', 'H': '\U0000210B', 'I': '\U00002110', 'J': '\U0001D4A5',
		'K': '\U0001D4A6', 'L': '\U00002112', 'M': '\U00002133', 'N': '\U0001D4A9', 'O': '\U0001D4AA',
		'P': '\U0001D4AB', 'Q': '\U0001D4AC', 'R': '\U0000211B', 'S': '\U0001D4AE', 'T': '\U0001D4AF',
		'U': '\U0001D4B0', 'V': '\U0001D4B1', 'W': '\U0001D4B2', 'X': '\U0001D4B3', 'Y': '\U0001D4B4',
		'Z': '\U0001D4B5',

		'a': '\U0001D4B6', 'b': '\U0001D4B7', 'c': '\U0001D4B8', 'd': '\U0001D4B9', 'e': '\U0000212F',
		'f': '\U0001D4BB', 'g': '\U0000210A', 'h': '\U0001D4BD', 'i': '\U0001D4BE', 'j': '\U0001D4BF',
		'k': '\U0001D4C0', 'l': '\U0001D4C1', 'm': '\U0001D4C2', 'n': '\U0001D4C3', 'o': '\U00002134',
		'p': '\U0001D4C5', 'q': '\U0001D4C6', 'r': '\U0001D4C7', 's': '\U0001D4C8', 't': '\U0001D4C9',
		'u': '\U0001D4CA', 'v': '\U0001D4CB', 'w': '\U0001D4CC', 'x': '\U0001D4CD', 'y': '\U0001D4CE',
		'z': '\U0001D4CF',
	}

	return translateMap(text, m)
}

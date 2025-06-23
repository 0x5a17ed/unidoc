package unidoc

// toBoldSansSerifText converts regular text to a bold sans-serif style using mathematical Unicode characters.
func toBoldSansSerifText(text string) string {
	// Mathematical Sans-Serif Bold Unicode mapping
	m := map[rune]rune{
		'A': '\U0001D5D4', 'B': '\U0001D5D5', 'C': '\U0001D5D6', 'D': '\U0001D5D7', 'E': '\U0001D5D8',
		'F': '\U0001D5D9', 'G': '\U0001D5DA', 'H': '\U0001D5DB', 'I': '\U0001D5DC', 'J': '\U0001D5DD',
		'K': '\U0001D5DE', 'L': '\U0001D5DF', 'M': '\U0001D5E0', 'N': '\U0001D5E1', 'O': '\U0001D5E2',
		'P': '\U0001D5E3', 'Q': '\U0001D5E4', 'R': '\U0001D5E5', 'S': '\U0001D5E6', 'T': '\U0001D5E7',
		'U': '\U0001D5E8', 'V': '\U0001D5E9', 'W': '\U0001D5EA', 'X': '\U0001D5EB', 'Y': '\U0001D5EC',
		'Z': '\U0001D5ED',

		'a': '\U0001D5EE', 'b': '\U0001D5EF', 'c': '\U0001D5F0', 'd': '\U0001D5F1', 'e': '\U0001D5F2',
		'f': '\U0001D5F3', 'g': '\U0001D5F4', 'h': '\U0001D5F5', 'i': '\U0001D5F6', 'j': '\U0001D5F7',
		'k': '\U0001D5F8', 'l': '\U0001D5F9', 'm': '\U0001D5FA', 'n': '\U0001D5FB', 'o': '\U0001D5FC',
		'p': '\U0001D5FD', 'q': '\U0001D5FE', 'r': '\U0001D5FF', 's': '\U0001D600', 't': '\U0001D601',
		'u': '\U0001D602', 'v': '\U0001D603', 'w': '\U0001D604', 'x': '\U0001D605', 'y': '\U0001D606',
		'z': '\U0001D607',

		'0': '\U0001D7EC', '1': '\U0001D7ED', '2': '\U0001D7EE', '3': '\U0001D7EF', '4': '\U0001D7F0',
		'5': '\U0001D7F1', '6': '\U0001D7F2', '7': '\U0001D7F3', '8': '\U0001D7F4', '9': '\U0001D7F5',
	}

	return translateMap(text, m)
}

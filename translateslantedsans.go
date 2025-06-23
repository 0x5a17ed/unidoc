package unidoc

// toSlantedSansSerifText converts regular text to mathematical sans-serif italic Unicode (𝘛𝘩𝘪𝘴)
func toSlantedSansSerifText(text string) string {
	// Mathematical Sans-Serif Italic Unicode mapping
	m := map[rune]rune{
		'A': '\U0001D608', 'B': '\U0001D609', 'C': '\U0001D60A', 'D': '\U0001D60B', 'E': '\U0001D60C',
		'F': '\U0001D60D', 'G': '\U0001D60E', 'H': '\U0001D60F', 'I': '\U0001D610', 'J': '\U0001D611',
		'K': '\U0001D612', 'L': '\U0001D613', 'M': '\U0001D614', 'N': '\U0001D615', 'O': '\U0001D616',
		'P': '\U0001D617', 'Q': '\U0001D618', 'R': '\U0001D619', 'S': '\U0001D61A', 'T': '\U0001D61B',
		'U': '\U0001D61C', 'V': '\U0001D61D', 'W': '\U0001D61E', 'X': '\U0001D61F', 'Y': '\U0001D620',
		'Z': '\U0001D621',

		'a': '\U0001D622', 'b': '\U0001D623', 'c': '\U0001D624', 'd': '\U0001D625', 'e': '\U0001D626',
		'f': '\U0001D627', 'g': '\U0001D628', 'h': '\U0001D629', 'i': '\U0001D62A', 'j': '\U0001D62B',
		'k': '\U0001D62C', 'l': '\U0001D62D', 'm': '\U0001D62E', 'n': '\U0001D62F', 'o': '\U0001D630',
		'p': '\U0001D631', 'q': '\U0001D632', 'r': '\U0001D633', 's': '\U0001D634', 't': '\U0001D635',
		'u': '\U0001D636', 'v': '\U0001D637', 'w': '\U0001D638', 'x': '\U0001D639', 'y': '\U0001D63A',
		'z': '\U0001D63B',
	}

	return translateMap(text, m)
}

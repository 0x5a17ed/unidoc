package unidoc

// Config holds the configuration for the Unicode renderer.
type Config struct {
	ItalicStyle ItalicStyle // Style for italic text: "markers", "script", "sans-italic"
	StrongStyle StrongStyle // Style for strong text: "plain", "markers", "math"
}

// DefaultConfig returns the default configuration for the Unicode renderer.
func DefaultConfig() Config {
	return Config{
		ItalicStyle: ItalicStyleSlantedSansSerif, // Default italic style
		StrongStyle: StrongStyleBoldSansSerif,    // Default strong style
	}
}

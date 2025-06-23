package unidoc

import (
	"fmt"
)

type ItalicStyle int

const (
	ItalicStylePlain            ItalicStyle = iota // Use plain style for italics
	ItalicStyleMarkers                             // Use simple markers for italics
	ItalicStyleScript                              // Use script style for italics
	ItalicStyleSlantedSansSerif                    // Use slanted sans-serif italic style for italics
)

// UnmarshalText implements the encoding.TextUnmarshaler interface for ItalicStyle.
func (s *ItalicStyle) UnmarshalText(text []byte) error {
	switch string(text) {
	case "plain":
		*s = ItalicStylePlain
	case "markers":
		*s = ItalicStyleMarkers
	case "script":
		*s = ItalicStyleScript
	case "slanted-sans-serif":
		*s = ItalicStyleSlantedSansSerif
	default:
		return fmt.Errorf("invalid italic style: %s", text)
	}
	return nil
}

// Set implements the pflag.Value interface for ItalicStyle.
func (s *ItalicStyle) Set(value string) error {
	return s.UnmarshalText([]byte(value))
}

// String implements the pflag.Value interface for ItalicStyle.
func (s *ItalicStyle) String() string {
	switch *s {
	case ItalicStylePlain:
		return "plain"
	case ItalicStyleMarkers:
		return "markers"
	case ItalicStyleScript:
		return "script"
	case ItalicStyleSlantedSansSerif:
		return "slanted-sans-serif"
	default:
		return "unknown"
	}
}

// Type implements the pflag.Value interface for ItalicStyle.
func (s *ItalicStyle) Type() string {
	return "italicStyle"
}

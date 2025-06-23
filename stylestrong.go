package unidoc

import (
	"fmt"
	"strings"
)

type StrongStyle int

const (
	StrongStylePlain         StrongStyle = iota // Use plain style for strong text
	StrongStyleMarkers                          // Use bold style for strong text
	StrongStyleBoldSansSerif                    // Use mathematical bold sans-serif for strong text
)

// UnmarshalText implements the encoding.TextUnmarshaler interface for StrongStyle.
func (s *StrongStyle) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "plain":
		*s = StrongStylePlain
	case "markers":
		*s = StrongStyleMarkers
	case "bold-sans-serif":
		*s = StrongStyleBoldSansSerif
	default:
		return fmt.Errorf("invalid strong style: %s", text)
	}
	return nil
}

// Set implements the pflag.Value interface for StrongStyle.
func (s *StrongStyle) Set(value string) error {
	return s.UnmarshalText([]byte(value))
}

// String implements the pflag.Value interface for StrongStyle.
func (s *StrongStyle) String() string {
	switch *s {
	case StrongStylePlain:
		return "plain"
	case StrongStyleMarkers:
		return "markers"
	case StrongStyleBoldSansSerif:
		return "bold-sans-serif"
	default:
		return "unknown"
	}
}

// Type implements the pflag.Value interface for StrongStyle.
func (s *StrongStyle) Type() string {
	return "strongStyle"
}

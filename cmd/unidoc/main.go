package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"

	"github.com/0x5a17ed/unidoc"
)

func showHelp() {
	_, _ = fmt.Fprintf(os.Stderr, `
unidoc - render markdown with Unicode symbols

Renders Markdown documents with Unicode symbols for emphasis, headers, 
lists, and more.

Usage:
  unidoc [OPTION]... [FILE]

Options:
`)

	pflag.PrintDefaults()

	_, _ = fmt.Fprintf(os.Stderr, `
Italic Styles:
  plain                       use regular text, no special formatting
  markers                     use simple markers around italic text: *text*
  script                      use Unicode script letters: 𝒯𝒽𝒾𝓈 𝒾𝓈 𝒾𝓉𝒶𝓁𝒾𝒸
  slanted-sans-serif          use Unicode sans-serif italic: 𝘛𝘩𝘪𝘴 𝘪𝘴 𝘪𝘵𝘢𝘭𝘪𝘤

Strong Styles:
  plain                       use regular text, no special formatting
  markers                     use simple markers around strong text: **text**
  bold-sans-serif             use mathematical bold sans-serif: 𝗧𝗵𝗶𝘀 𝗶𝘀 𝗯𝗼𝗹𝗱

Examples:
  echo '# Hello World' | unidoc
  unidoc README.md
  unidoc --italic script < document.md
`)
}

func mainE() error {
	config := unidoc.DefaultConfig()

	// Define flags
	var (
		showHelpFlag = pflag.BoolP("help", "h", false, "Show help information")
	)
	pflag.Var(&config.ItalicStyle, "italic-style", "style for italic text")
	pflag.Var(&config.StrongStyle, "strong-style", "style for strong text")

	pflag.Usage = showHelp
	pflag.Parse()

	// Show help if requested
	if *showHelpFlag {
		showHelp()
		return nil
	}

	var (
		content  []byte
		err      error
		fileName = pflag.Arg(0) // Get the first positional argument as file name
	)

	// Read from file if specified.
	if fileName != "" && fileName != "-" {
		content, err = os.ReadFile(fileName)
	} else {
		// Read from stdin.
		fileName = "stdin"
		content, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", fileName, err)
	}

	content = bytes.TrimSpace(content)
	if len(content) > 0 {
		result, err := unidoc.Convert(content, config)
		if err != nil {
			return fmt.Errorf("conversion error: %w", err)
		}
		fmt.Println(result)
	}
	return nil
}

func main() {
	if err := mainE(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

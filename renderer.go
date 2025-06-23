package unidoc

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// UnicodeRenderer implements a custom renderer for pure Unicode text output
type UnicodeRenderer struct {
	config Config

	listLevel       int
	blockquoteLevel int
	inHeader        bool
	inStrong        bool
	inItalic        bool
	inListItem      bool
	inBlockquote    bool
	listNumbers     []int  // Stack to track current numbers for nested ordered lists
	isOrdered       []bool // Stack to track if current lists are ordered
}

// NewUnicodeRenderer creates a new Unicode text renderer
func NewUnicodeRenderer(config Config) *UnicodeRenderer {
	return &UnicodeRenderer{
		config: config,
	}
}

// RegisterFuncs registers the renderer for all node types
func (r *UnicodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Block nodes
	reg.Register(gast.KindDocument, r.renderDocument)
	reg.Register(gast.KindHeading, r.renderHeading)
	reg.Register(gast.KindParagraph, r.renderParagraph)
	reg.Register(gast.KindList, r.renderList)
	reg.Register(gast.KindListItem, r.renderListItem)
	reg.Register(gast.KindBlockquote, r.renderBlockquote)
	reg.Register(gast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(gast.KindFencedCodeBlock, r.renderCodeBlock)
	reg.Register(gast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(gast.KindThematicBreak, r.renderThematicBreak)

	// Inline nodes
	reg.Register(gast.KindText, r.renderText)
	reg.Register(gast.KindEmphasis, r.renderEmphasis)
	reg.Register(gast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(gast.KindLink, r.renderLink)
	reg.Register(gast.KindImage, r.renderImage)
	reg.Register(gast.KindAutoLink, r.renderAutoLink)
	reg.Register(gast.KindRawHTML, r.renderRawHTML)
	reg.Register(gast.KindTextBlock, r.renderTextBlock)
}

// Document renderer
func (r *UnicodeRenderer) renderDocument(
	_ util.BufWriter,
	_ []byte,
	_ gast.Node,
	_ bool,
) (gast.WalkStatus, error) {
	return gast.WalkContinue, nil
}

// Heading renderer
func (r *UnicodeRenderer) renderHeading(
	w util.BufWriter,
	_ []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	n := node.(*gast.Heading)
	if entering {
		r.inHeader = true
		// Render headers with Unicode box drawing characters and emphasis
		level := n.Level
		var prefix string
		switch level {
		case 1:
			prefix = "█" // Full block for H1
		case 2:
			prefix = "▓" // Medium shade for H2
		case 3:
			prefix = "▒" // Light shade for H3
		case 4:
			prefix = "░" // Very light shade for H4
		case 5:
			prefix = "▫" // Small square for H5
		default:
			prefix = "▪" // Small black square for H6+
		}

		if level > 1 && level < 5 {
			prefix = strings.Repeat(prefix, level)
		}

		if _, err := w.WriteString(prefix + " "); err != nil {
			return gast.WalkStop, err
		}
	} else {
		r.inHeader = false
		// Add underline for H1 and H2
		if n.Level <= 2 {
			if _, err := w.WriteString("\n"); err != nil {
				return gast.WalkStop, err
			}
			char := "═"
			if n.Level == 2 {
				char = "─"
			}
			// Estimate header length (rough approximation)
			underline := strings.Repeat(char, 50)
			if _, err := w.WriteString(underline); err != nil {
				return gast.WalkStop, err
			}
		}
		if _, err := w.WriteString("\n\n"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

// Paragraph renderer
func (r *UnicodeRenderer) renderParagraph(
	w util.BufWriter,
	_ []byte,
	_ gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if !entering {
		// Only add extra newlines if we're not inside a list item
		if !r.inListItem {
			if _, err := w.WriteString("\n\n"); err != nil {
				return gast.WalkStop, err
			}
		}
	}
	return gast.WalkContinue, nil
}

// Text renderer
func (r *UnicodeRenderer) renderText(
	w util.BufWriter,
	source []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}

	n := node.(*gast.Text)
	segment := n.Segment
	text := string(segment.Value(source))

	text = r.toSmartDashes(text)

	switch {
	case r.inHeader:
		text = toBoldSansSerifText(text)
	case r.inStrong:
		switch r.config.StrongStyle {
		case StrongStylePlain:
			// Use plain strong style, no conversion needed
		case StrongStyleMarkers:
			// Use simple markers for strong text
			text = fmt.Sprintf("**%s**", text)
		case StrongStyleBoldSansSerif:
			// Use mathematical bold sans-serif for strong text
			text = toBoldSansSerifText(text)
		}

	case r.inItalic:
		switch r.config.ItalicStyle {
		case ItalicStylePlain:
			// Use plain italic style, no conversion needed
		case ItalicStyleMarkers:
			// Use simple markers for italic text
			text = fmt.Sprintf("*%s*", text)
		case ItalicStyleScript:
			text = toItalicScriptText(text)
		case ItalicStyleSlantedSansSerif:
			text = toSlantedSansSerifText(text)
		}
	}

	if _, err := w.WriteString(text); err != nil {
		return gast.WalkStop, err
	}

	// Check if this text node ends with a hard line break (double space in markdown)
	if n.HardLineBreak() {
		if _, err := w.WriteString("\n"); err != nil {
			return gast.WalkStop, err
		}
		// Add appropriate indentation based on context
		if r.inBlockquote {
			// Add blockquote prefixes for continuation lines
			prefix := strings.Repeat("┃ ", r.blockquoteLevel)
			if _, err := w.WriteString(prefix); err != nil {
				return gast.WalkStop, err
			}
		} else if r.inListItem {
			// Add list indentation for continuation lines
			indent := strings.Repeat("  ", r.listLevel)
			if _, err := w.WriteString(indent); err != nil {
				return gast.WalkStop, err
			}
		}
	}

	return gast.WalkContinue, nil
}

// Emphasis renderer (handles both bold and italic)
func (r *UnicodeRenderer) renderEmphasis(
	w util.BufWriter,
	_ []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	n := node.(*gast.Emphasis)

	if n.Level > 1 { // Strong/Bold (**text**)
		r.inStrong = entering

	} else { // Italic (*text*)
		r.inItalic = entering
	}

	return gast.WalkContinue, nil
}

// List renderer
func (r *UnicodeRenderer) renderList(
	w util.BufWriter,
	_ []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	n := node.(*gast.List)

	if entering {
		// Add newline before nested lists (but not for top-level lists)
		if r.listLevel > 0 {
			if _, err := w.WriteString("\n"); err != nil {
				return gast.WalkStop, err
			}
		}

		r.listLevel++
		// Check if this is an ordered list
		isOrdered := n.IsOrdered()
		r.isOrdered = append(r.isOrdered, isOrdered)

		// Initialize counter for ordered lists
		if isOrdered {
			start := n.Start
			if start == 0 {
				start = 1
			}
			r.listNumbers = append(r.listNumbers, start)
		} else {
			r.listNumbers = append(r.listNumbers, 0) // 0 for unordered
		}
	} else {
		r.listLevel--
		// Pop from stacks
		if len(r.isOrdered) > 0 {
			r.isOrdered = r.isOrdered[:len(r.isOrdered)-1]
		}
		if len(r.listNumbers) > 0 {
			r.listNumbers = r.listNumbers[:len(r.listNumbers)-1]
		}
		if r.listLevel == 0 {
			if _, err := w.WriteString("\n"); err != nil {
				return gast.WalkStop, err
			}
		}
	}
	return gast.WalkContinue, nil
}

// ListItem renderer
func (r *UnicodeRenderer) renderListItem(
	w util.BufWriter,
	_ []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		r.inListItem = true
		indent := strings.Repeat("  ", r.listLevel-1)

		// Determine if current list is ordered
		isCurrentOrdered := len(r.isOrdered) > 0 && r.isOrdered[len(r.isOrdered)-1]

		if isCurrentOrdered {
			// Get current number and increment for next item
			currentNum := r.listNumbers[len(r.listNumbers)-1]
			r.listNumbers[len(r.listNumbers)-1]++

			// Use fancy Unicode numbering based on nesting level
			marker := r.getOrderedMarker(currentNum, r.listLevel)
			if _, err := w.WriteString(fmt.Sprintf("%s%s ", indent, marker)); err != nil {
				return gast.WalkStop, err
			}
		} else {
			// Use Unicode bullets for unordered lists
			bullets := []string{"•", "◦", "▪", "▫", "‣", "⁃"}
			bullet := bullets[(r.listLevel-1)%len(bullets)]
			if _, err := w.WriteString(fmt.Sprintf("%s%s ", indent, bullet)); err != nil {
				return gast.WalkStop, err
			}
		}
	} else {
		r.inListItem = false

		// Check if parent list is loose (has blank lines between items)
		// by checking if this list item contains paragraph nodes
		hasParaChild := false
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if child.Kind() == gast.KindParagraph {
				hasParaChild = true
				break
			}
		}

		// Add extra spacing for loose lists
		if hasParaChild {
			if _, err := w.WriteString("\n\n"); err != nil {
				return gast.WalkStop, err
			}
		} else {
			if _, err := w.WriteString("\n"); err != nil {
				return gast.WalkStop, err
			}
		}
	}
	return gast.WalkContinue, nil
}

// Blockquote renderer
func (r *UnicodeRenderer) renderBlockquote(
	w util.BufWriter,
	_ []byte,
	_ gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		r.blockquoteLevel++
		// Render the appropriate number of ┃ symbols for nesting
		prefix := strings.Repeat("┃ ", r.blockquoteLevel)
		if _, err := w.WriteString(prefix); err != nil {
			return gast.WalkStop, err
		}
	} else {
		r.blockquoteLevel--
		if _, err := w.WriteString("\n\n"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

// CodeBlock renderer
func (r *UnicodeRenderer) renderCodeBlock(
	w util.BufWriter,
	source []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}

	// Extract code content differently for fenced vs regular code blocks
	var code string

	if fcb, ok := node.(*gast.FencedCodeBlock); ok {
		// Fenced code block
		var buf bytes.Buffer
		for i := 0; i < fcb.Lines().Len(); i++ {
			line := fcb.Lines().At(i)
			buf.Write(line.Value(source))
		}
		code = buf.String()
	} else {
		// Regular code block
		var buf bytes.Buffer
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if text, ok := child.(*gast.Text); ok {
				segment := text.Segment
				buf.Write(segment.Value(source))
			}
		}
		code = buf.String()
	}

	lines := strings.Split(strings.TrimRight(code, "\n"), "\n")

	// Top border
	if _, err := w.WriteString("┌" + strings.Repeat("─", 64) + "┐\n"); err != nil {
		return gast.WalkStop, err
	}

	// Code content
	for _, line := range lines {
		// Expand tabs to spaces (4 spaces per tab)
		line = strings.ReplaceAll(line, "\t", "    ")

		// Ensure line doesn't exceed box width
		if utf8.RuneCountInString(line) > 58 {
			line = line[:58] + "…"
		}
		padding := 64 - utf8.RuneCountInString(line)
		if _, err := w.WriteString("│ " + line + strings.Repeat(" ", padding-2) + " │\n"); err != nil {
			return gast.WalkStop, err
		}
	}

	// Bottom border
	if _, err := w.WriteString("└" + strings.Repeat("─", 64) + "┘\n\n"); err != nil {
		return gast.WalkStop, err
	}

	return gast.WalkSkipChildren, nil
}

// CodeSpan renderer (inline code)
func (r *UnicodeRenderer) renderCodeSpan(
	w util.BufWriter,
	_ []byte,
	_ gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("⌜"); err != nil {
			return gast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("⌝"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

// Link renderer
func (r *UnicodeRenderer) renderLink(
	w util.BufWriter,
	_ []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		// Use Unicode arrow for links
		if _, err := w.WriteString("["); err != nil {
			return gast.WalkContinue, err
		}
	} else {
		n := node.(*gast.Link)
		url := string(n.Destination)
		if _, err := w.WriteString(fmt.Sprintf("] 🔗 <%s>", url)); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

// Image renderer
func (r *UnicodeRenderer) renderImage(
	w util.BufWriter,
	_ []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		n := node.(*gast.Image)
		alt := string(n.Title)
		if alt == "" {
			alt = "Image"
		}
		url := string(n.Destination)
		if _, err := w.WriteString(fmt.Sprintf("🖼️  %s <%s>", alt, url)); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkSkipChildren, nil
}

// ThematicBreak renderer (horizontal rule)
func (r *UnicodeRenderer) renderThematicBreak(
	w util.BufWriter,
	_ []byte,
	_ gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("\n" + strings.Repeat("═", 60) + "\n\n"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

// Fallback renderers for other node types
func (r *UnicodeRenderer) renderAutoLink(
	w util.BufWriter,
	source []byte,
	node gast.Node,
	entering bool,
) (gast.WalkStatus, error) {
	return r.renderLink(w, source, node, entering)
}

func (r *UnicodeRenderer) renderRawHTML(
	_ util.BufWriter,
	_ []byte,
	_ gast.Node,
	_ bool,
) (gast.WalkStatus, error) {
	// Skip HTML content
	return gast.WalkSkipChildren, nil
}

func (r *UnicodeRenderer) renderHTMLBlock(
	_ util.BufWriter,
	_ []byte,
	_ gast.Node,
	_ bool,
) (gast.WalkStatus, error) {
	// Skip HTML blocks
	return gast.WalkSkipChildren, nil
}

func (r *UnicodeRenderer) renderTextBlock(
	_ util.BufWriter,
	_ []byte,
	_ gast.Node,
	_ bool,
) (gast.WalkStatus, error) {
	return gast.WalkContinue, nil
}

// getOrderedMarker returns fancy Unicode numbering based on the number and nesting level
func (r *UnicodeRenderer) getOrderedMarker(num int, level int) string {
	switch level {
	case 1:
		// Level 1: Circled numbers ① ② ③ etc.
		return r.getCircledNumber(num)
	case 2:
		// Level 2: Parenthesized numbers ⑴ ⑵ ⑶ etc.
		return r.getParenthesizedNumber(num)
	case 3:
		// Level 3: Circled letters 🅐 🅑 🅒 etc.
		return r.getCircledLetter(num)
	case 4:
		// Level 4: Roman numerals ⅰ ⅱ ⅲ etc.
		return r.getRomanNumeral(num, false) // lowercase
	case 5:
		// Level 5: Roman numerals uppercase Ⅰ Ⅱ Ⅲ etc.
		return r.getRomanNumeral(num, true) // uppercase
	default:
		// Fallback: regular numbers with fancy formatting
		return fmt.Sprintf("⟨%d⟩", num)
	}
}

// getCircledNumber returns circled Unicode numbers ① ② ③ etc.
func (r *UnicodeRenderer) getCircledNumber(num int) string {
	// Unicode circled numbers 1-20
	circledNums := []string{
		"①", "②", "③", "④", "⑤", "⑥", "⑦", "⑧", "⑨", "⑩",
		"⑪", "⑫", "⑬", "⑭", "⑮", "⑯", "⑰", "⑱", "⑲", "⑳",
	}

	if num >= 1 && num <= 20 {
		return circledNums[num-1]
	}
	// Fallback for numbers > 20
	return fmt.Sprintf("⓪%d", num)
}

// getParenthesizedNumber returns parenthesized Unicode numbers ⑴ ⑵ ⑶ etc.
func (r *UnicodeRenderer) getParenthesizedNumber(num int) string {
	// Unicode parenthesized numbers 1-20
	parenNums := []string{
		"⑴", "⑵", "⑶", "⑷", "⑸", "⑹", "⑺", "⑻", "⑼", "⑽",
		"⑾", "⑿", "⒀", "⒁", "⒂", "⒃", "⒄", "⒅", "⒆", "⒇",
	}

	if num >= 1 && num <= 20 {
		return parenNums[num-1]
	}
	// Fallback for numbers > 20
	return fmt.Sprintf("(%d)", num)
}

// getCircledLetter returns circled Unicode letters 🅐 🅑 🅒 etc.
func (r *UnicodeRenderer) getCircledLetter(num int) string {
	// Unicode circled letters A-Z
	circledLetters := []string{
		"🅐", "🅑", "🅒", "🅓", "🅔", "🅕", "🅖", "🅗", "🅘", "🅙",
		"🅚", "🅛", "🅜", "🅝", "🅞", "🅟", "🅠", "🅡", "🅢", "🅣",
		"🅤", "🅥", "🅦", "🅧", "🅨", "🅩",
	}

	if num >= 1 && num <= 26 {
		return circledLetters[num-1]
	}
	// Cycle through letters for numbers > 26
	letter := circledLetters[(num-1)%26]
	cycle := (num-1)/26 + 1
	if cycle > 1 {
		return fmt.Sprintf("%s·%d", letter, cycle)
	}
	return letter
}

// getRomanNumeral returns Roman numerals in Unicode
func (r *UnicodeRenderer) getRomanNumeral(num int, uppercase bool) string {
	if num <= 0 || num > 50 {
		// Fallback for out of range
		if uppercase {
			return fmt.Sprintf("Ⅻ·%d", num)
		}
		return fmt.Sprintf("ⅻ·%d", num)
	}

	var lowerRoman = []string{
		"ⅰ", "ⅱ", "ⅲ", "ⅳ", "ⅴ", "ⅵ", "ⅶ", "ⅷ", "ⅸ", "ⅹ",
		"ⅺ", "ⅻ",
	}

	var upperRoman = []string{
		"Ⅰ", "Ⅱ", "Ⅲ", "Ⅳ", "Ⅴ", "Ⅵ", "Ⅶ", "Ⅷ", "Ⅸ", "Ⅹ",
		"Ⅺ", "Ⅻ",
	}

	// For numbers 1-12, use dedicated Unicode Roman numerals
	if num <= 12 {
		if uppercase {
			return upperRoman[num-1]
		}
		return lowerRoman[num-1]
	}

	// For numbers 13+, construct Roman numerals
	// This is a simplified version - full Roman numeral conversion would be more complex
	if uppercase {
		return fmt.Sprintf("Ⅻ+%d", num-12)
	}
	return fmt.Sprintf("ⅻ+%d", num-12)
}

func (r *UnicodeRenderer) toSmartDashes(text string) string {
	// Convert in order from longest to shortest to avoid conflicts
	// --- → em dash (—)
	text = strings.ReplaceAll(text, "---", "—")
	// -- → en dash (–)
	text = strings.ReplaceAll(text, "--", "–")

	return text
}

// Convert converts Markdown text to Unicode-rendered text
func Convert(inp []byte, config Config) (string, error) {
	// Create goldmark instance with Unicode renderer
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(NewUnicodeRenderer(config), 1000),
				),
			),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(inp, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}

	// Clean up multiple consecutive newlines
	result := buf.String()
	re := regexp.MustCompile(`\n{3,}`)
	result = re.ReplaceAllString(result, "\n\n")

	// Clean up multiple consecutive newlines
	return strings.TrimSpace(result), nil
}

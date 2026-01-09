package feature

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// CharacterAnalyzer analyzes text encoding and character properties
type CharacterAnalyzer struct{}

func NewCharacterAnalyzer() *CharacterAnalyzer {
	return &CharacterAnalyzer{}
}

func (ca *CharacterAnalyzer) ID() string {
	return "character"
}

func (ca *CharacterAnalyzer) Name() string {
	return "Character Analyzer"
}

func (ca *CharacterAnalyzer) Description() string {
	return "Analyze text encoding (UTF-8, UTF-16, ASCII, hex, binary)"
}

func (ca *CharacterAnalyzer) Help() string {
	return `Character Analyzer examines text and provides detailed encoding information:

• Decimal (ASCII/Unicode code points)
• Hexadecimal representation
• Binary representation
• UTF-8, UTF-16, UTF-32 byte counts
• Total character (rune) count

This is useful for understanding how text is encoded, debugging encoding issues,
or learning about character representations.`
}

func (ca *CharacterAnalyzer) Examples() []Example {
	return []Example{
		{Input: "Hello", Description: "Analyze simple ASCII text"},
		{Input: "你好", Description: "Analyze Chinese characters (multi-byte UTF-8)"},
		{Input: "🚀", Description: "Analyze emoji (4-byte UTF-8)"},
		{Input: "Café", Description: "Analyze text with accented characters"},
	}
}

func (ca *CharacterAnalyzer) Execute(input string) (string, error) {
	if input == "" {
		return "Please provide some text to analyze", nil
	}

	var result strings.Builder

	runeCount := utf8.RuneCountInString(input)
	utf8Bytes := len(input)
	utf16Bytes := calculateUTF16Bytes(input)
	utf32Bytes := runeCount * 4

	result.WriteString(fmt.Sprintf("Total Runes: %d characters\n", runeCount))
	result.WriteString(fmt.Sprintf("UTF-8:  %d bytes\n", utf8Bytes))
	result.WriteString(fmt.Sprintf("UTF-16: %d bytes\n", utf16Bytes))
	result.WriteString(fmt.Sprintf("UTF-32: %d bytes\n\n", utf32Bytes))

	result.WriteString("Decimal:     " + generateDecimal(input) + "\n")
	result.WriteString("Hexadecimal: " + generateHex(input) + "\n")
	result.WriteString("Binary:      " + truncate(generateBinary(input), 60) + "\n")

	return result.String(), nil
}

func calculateUTF16Bytes(text string) int {
	count := 0
	for _, r := range text {
		if r <= 0xFFFF {
			count += 2
		} else {
			count += 4
		}
	}
	return count
}

func generateDecimal(text string) string {
	var b strings.Builder
	for i, r := range text {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%d", r))
	}
	return b.String()
}

func generateHex(text string) string {
	var b strings.Builder
	for i, r := range text {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%X", r))
	}
	return b.String()
}

func generateBinary(text string) string {
	var b strings.Builder
	for i, r := range text {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%b", r))
	}
	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

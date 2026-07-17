package bot

import (
	"strings"
	"unicode"
)

func isNotLetter(r rune) bool {
	return !unicode.IsLetter(r)
}

func CleanAndSplit(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.FieldsFunc(lowered, isNotLetter)
	return words
}

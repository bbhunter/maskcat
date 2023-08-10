// Package utils contains functions for the main program
//
// The package structure is broken into two components:
//
// utils.go which contains the primary logic
// utils_test.go which contains unit tests
package utils

import (
	"regexp"
	"strings"

	"github.com/jakewnuk/maskcat/pkg/models"
)

// ConstructReplacements create an array mapping which characters to replace
//
// This function accepts the characters "ulds" in order to generate a map
// - u for uppercase characters
// - l for lowercase characters
// - d for numerical characters
// - s for special characters
//
// Args:
//
//	str (string): Input string
//
// Returns:
//
//	args ([]string): Map of replacement characters
func ConstructReplacements(str string) []string {
	var lowerArgs, upperArgs, digitArgs, args []string
	for c := 'a'; c <= 'z'; c++ {
		lowerArgs = append(lowerArgs, string(c), "?l")
	}
	for c := 'A'; c <= 'Z'; c++ {
		upperArgs = append(upperArgs, string(c), "?u")
	}
	for c := '0'; c <= '9'; c++ {
		digitArgs = append(digitArgs, string(c), "?d")
	}
	specialChars := " !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~'"
	specialArgs := make([]string, len(specialChars)*2)
	for i, c := range specialChars {
		specialArgs[i*2] = string(c)
		specialArgs[i*2+1] = "?s"
	}

	if strings.Contains(str, "l") {
		args = append(args, lowerArgs...)
	}

	if strings.Contains(str, "u") {
		args = append(args, upperArgs...)
	}

	if strings.Contains(str, "d") {
		args = append(args, digitArgs...)
	}

	if strings.Contains(str, "s") {
		args = append(args, specialArgs...)
	}

	return args
}

// MakeMask performs substitution to make masks
//
// Args:
//
//	str (string): String to turn into a mask
//	replacements ([]string): Map of which characters to replace
//
// Returns:
//
//	(string): Replaced string as a mask
func MakeMask(str string, replacements []string) string {
	return strings.NewReplacer(replacements...).Replace(str)
}

// MakeToken replaces all non-alpha characters to generate string tokens
//
// Args:
//
//	str (string): Input string
//
// Returns:
//
//	(string): String with only alpha characters
func MakeToken(str string) string {
	re := regexp.MustCompile(`[^a-zA-Z]+`)
	return re.ReplaceAllString(str, "")
}

// RemoveMaskCharacters will replace mask characters in a string with nothing
//
// Args:
//
//	str (string): Input string to replace
//
// Returns:
//
//	(string): String with replaced characters
func RemoveMaskCharacters(str string) string {
	return strings.NewReplacer("?u", "", "?l", "", "?d", "", "?b", "", "?s", "").Replace(str)
}

// TestComplexity tests the complexity of an input mask
//
// Args:
//
//	str (string): Input string to test
//
// Returns:
//
//	(int): Complexity score as an integer
func TestComplexity(str string) int {
	complexity := 0
	charTypes := []string{"?u", "?l", "?d", "?s", "?b"}
	for _, charType := range charTypes {
		if strings.Contains(str, charType) {
			complexity++
		}
	}
	return complexity
}

// TestEntropy calculates mask entropy of an input mask
//
// Args:
//
//	str (string): Input string to test
//
// Returns:
//
//	(int): Entropy score as an integer
func TestEntropy(str string) int {
	entropy := 0
	charTypes := []struct {
		charType string
		count    int
	}{
		{"?u", 26},
		{"?l", 26},
		{"?d", 10},
		{"?s", 33},
		{"?b", 256},
	}
	for _, ct := range charTypes {
		entropy += strings.Count(str, ct.charType) * ct.count
	}
	return entropy
}

// ChunkString splits string into chunks of a given size
//
// Args:
//
//	str (string): Input string to split
//	chunkSize (int): Chunk size to split
//
// Returns:
//
//	chunks ([]string): Map of chunked string
func ChunkString(str string, chunkSize int) []string {
	if len(str) == 0 {
		return nil
	}
	if chunkSize >= len(str) {
		return []string{str}
	}
	var chunks []string
	for i := 0; i < len(str); i += chunkSize {
		end := i + chunkSize
		if end > len(str) {
			end = len(str)
		}
		chunks = append(chunks, str[i:end])
	}
	return chunks
}

// ReplaceWord replaces a mask within an input string with a provided value
//
// Args:
//
//	word (string): Word to make replacements in
//	mask (string): Mask of the word
//	value (string): String to replace into the word
//	replacements ([]string): Replacement array used for the value parameter
//
// Returns:
//
//	newWord (string): Replaced word with value
func ReplaceWord(word string, mask string, value string, replacements []string) string {
	tokenmask := MakeMask(value, replacements)
	tokenmask = models.EnsureValidMask(tokenmask)

	if strings.Contains(mask, tokenmask) {
		newWord := strings.Replace(mask, tokenmask, value, 1)
		newWord = strings.NewReplacer("?u", "?", "?l", "?", "?b", "?", "?d", "?", "?s", "?").Replace(newWord)

		var builder strings.Builder
		builder.Grow(len(newWord))
		for i, r := range newWord {
			if r == '?' && i < len(word) {
				builder.WriteRune(rune(word[i]))
			} else {
				builder.WriteRune(r)
			}
		}
		newWord = builder.String()

		if strings.Contains(newWord, value) && newWord != value {
			return newWord
		}
	}
	return ""
}

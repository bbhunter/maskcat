// Package utils contains functions for the main program
//
// The package structure is broken into two components:
//
// utils.go which contains the primary logic
// utils_test.go which contains unit tests
package utils

import (
	"encoding/hex"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jakewnuk/maskcat/v2/pkg/models"
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

// MakeToken parses out tokens into an array
//   - Parses out camel case
//   - Parses out digit boundaries
//   - Parses out special char boundaries
//
// Args:
//
//	str (string): Input string
//
// Returns:
//
//	result ([]string): Tokens from input string
func MakeToken(str string) []string {
	re1 := regexp.MustCompile(`[A-Z][a-z]*|\d+|[^\dA-Z]+`)
	preArray := re1.FindAllString(str, -1)
	re2 := regexp.MustCompile(`[A-Z][a-z]*|\d+|\W+|\w+`)
	array := []string{}
	for _, s := range preArray {
		array = append(array, re2.FindAllString(s, -1)...)
	}
	re3 := regexp.MustCompile(`[^a-zA-Z]+`)
	array = append(array, re3.ReplaceAllString(str, ""))

	result := []string{}
	for _, word := range array {
		if word != " " {
			result = append(result, word)
		}
	}

	return result
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

// ReplaceAtIndex replaces a rune at index in string
//
// Args:
//
//	input (string): Input string to replace into
//	r (rune): Rune to replace into input
//	i (int): Position to replace at
//
// Returns:
//
//	out (string): Replaced string
func ReplaceAtIndex(input string, r rune, i int) string {
	if i < 0 || i >= len(input) {
		os.Exit(1)
	}
	out := []rune(input)
	if i >= 0 && i < len(out) {
		out[i] = r
		// In instances where i is out of bounds go to the end
	} else if i >= 0 && i == len(out) {
		out[len(out)-1] = r
	}
	return string(out)
}

// ReplaceWordByMask replaces a mask within an input string with a provided value
//
// Args:
//
//	word (string): Word to make replacements in
//	mask (string): Mask of the word
//	value (string): String to replace into the word
//	replacements ([]string): Replacement array used for the value parameter
//	fuzz (int): Amount of extra replacement characters to add
//
// Returns:
//
//	newWord (string): Replaced word with value
func ReplaceWordByMask(word string, mask string, value string, replacements []string, numOfReplacements int, fuzz int) string {
	tokenmask := MakeMask(value, replacements)
	tokenmask = models.EnsureValidMask(tokenmask)

	if fuzz > 0 {
		if fuzz > len(mask)/2 {
			fuzz = (len(mask) / 2)
		}
		mask = mask + mask[len(mask)-(fuzz*2):]
	}

	if strings.Contains(mask, tokenmask) {
		newword := strings.Replace(mask, tokenmask, value, numOfReplacements)
		newword = strings.NewReplacer("?u", "?", "?l", "?", "?b", "?", "?d", "?", "?s", "?").Replace(newword)

		for i := 0; i < len(word); {
			r, size := utf8.DecodeRuneInString(word[i:])

			if i < len(newword) && newword[i] == '?' {
				newword = ReplaceAtIndex(newword, r, i)
			}
			i += size
		}

		// NOTE: This introduces a known bug
		// If the first string contains "?" and a multibyte character the
		// output is malformed
		if !strings.Contains(word, "?") {
			newword = strings.ReplaceAll(newword, "?", "")
		}

		if strings.Contains(newword, value) && newword != value {
			return newword
		}
	}
	return ""
}

// DehexPlaintext decodes plaintext from $HEX[...] format
//
// Args:
//
//	s (string): The string to be dehexed
//
// Returns:
//
//	decoded (string): The decoded hex string
//	err (error): Error data
func DehexPlaintext(s string) (string, error) {
	s = strings.TrimPrefix(s, "$HEX[")
	s = strings.TrimSuffix(s, "]")
	decoded, err := hex.DecodeString(s)
	return string(decoded), err
}

// TestHexInput is used to identify plaintext in the $HEX[...] format
//
// Args:
//
//	s (str): The string to be evaluated
//
// Returns:
//
//	(bool): Returns true if it matches and false if it did not
func TestHexInput(s string) bool {
	var validateInput = regexp.MustCompile(`^\$HEX\[[a-zA-Z0-9]*\]$`).MatchString
	if validateInput(s) == false {
		return false
	}
	return true
}

// CreateRetainMask is used to create retain masks from input and a list of tokens
//
// # Retain masks are masks where keywords are prevented from being transformed
//
// Args:
//
//		stringWord (string): Input string to turn into a retain mask
//		retainTokens (map[string]struct{}): Tokens that should be kept
//	 args ([]string): Replacer arguments to use
//		doMultiByte	(bool): If the function should process multibyte text
//		doNumberOfReplacements (int): Number of tokens to keep in each string (default 1)
//
// Returns:
//
//	(string): Mask with any tokens retained
func CreateRetainMask(stringWord string, retainTokens map[string]struct{}, args []string, doMultiByte bool, doNumberOfReplacements int) string {
	// Create the retain mask
	result := []string{stringWord}

	// Iterate on tokens
	for value := range retainTokens {
		var temp []string

		// Iterate on item text
		for _, s := range result {
			split := strings.Split(s, value)

			// Iterate on exploded string
			for i, ss := range split {
				if ss != "" {
					temp = append(temp, ss)
				}

				// Check if its the last split string
				if i != len(split)-1 {
					temp = append(temp, value)
				}
			}
		}
		result = temp
	}

	kept := 0
	override := false
	forward := false

	for i, s := range result {
		if _, ok := retainTokens[s]; !ok || override {

			look := i + 1
			if look > len(result)-1 {
				look = len(result) - 1
			}

			// Limiting fowards to one for now
			if _, forwardOk := retainTokens[s+result[look]]; forwardOk && forward == false {
				forward = true
				continue
			} else if forward {
				forward = false
				kept++
				continue
			}

			mask := MakeMask(s, args)
			if doMultiByte {
				mask = models.EnsureValidMask(mask)
			}

			s = mask
			result[i] = s
		} else {
			if forward {
				forward = false
			}

			kept++
			if kept >= doNumberOfReplacements && doNumberOfReplacements != 0 {
				override = true
			}
		}
	}

	return strings.Join(result, "")

}

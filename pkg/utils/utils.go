// Package utils contains functions for the main maskcat program
package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jakewnuk/maskcat/pkg/models"
)

// ConstructReplacements create an array mapping which characters to replace
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

// MakeMask performs substitution to make HC masks
func MakeMask(str string, replacements []string) string {
	return strings.NewReplacer(replacements...).Replace(str)
}

// MakeToken replaces all non-alpha characters to generate tokens
func MakeToken(str string) string {
	re := regexp.MustCompile(`[^a-zA-Z]+`)
	return re.ReplaceAllString(str, "")
}

// MakePartialMask creates a partial Hashcat mask
func MakePartialMask(str string, replacements []string) string {
	return strings.NewReplacer(replacements...).Replace(str)
}

// RemoveMaskChars will replace mask characters in a string with nothing
func RemoveMaskChars(str string) string {
	return strings.NewReplacer("?u", "", "?l", "", "?d", "", "?b", "", "?s", "").Replace(str)
}

// TestComplexity tests the complexity of an input mask
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

// TestEntropy calculates mask entropy
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

// ChunkString splits string into chunks
func ChunkString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

// ReplaceWord replaces a mask within an input string with a value
func ReplaceWord(stringword, mask string, value string, replacements []string) string {
	tokenmask := MakeMask(value, replacements)
	tokenmask = models.ValidateMask(tokenmask)

	if strings.Contains(mask, tokenmask) {
		newword := strings.Replace(mask, tokenmask, value, -1)
		newword = strings.NewReplacer("?u", "?", "?l", "?", "?b", "?", "?d", "?", "?s", "?").Replace(newword)

		var builder strings.Builder
		builder.Grow(len(newword))
		for i, r := range newword {
			if r == '?' && i < len(stringword) {
				builder.WriteRune(rune(stringword[i]))
			} else {
				builder.WriteRune(r)
			}
		}
		newword = builder.String()

		if strings.Contains(newword, value) && newword != value {
			return newword
		}
	}
	return ""
}

// ConvertMultiByteString converts non-ascii characters to a valid format
func ConvertMultiByteString(str string) string {
	returnStr := ""
	for _, r := range str {
		if r > 127 {
			byteArr := []byte(string(r))
			for j := range byteArr {
				if j == len(byteArr)-1 {
					returnStr += fmt.Sprintf("?b")
				} else {
					returnStr += fmt.Sprintf("?b")
				}
			}
		} else {
			returnStr += fmt.Sprintf("%c", r)
		}
	}
	return returnStr
}

// CheckError is a general error handler
func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

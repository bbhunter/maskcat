// Package utils contains functions for the main maskcat program
package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

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

// RemoveDuplicateStr removes duplicate strings from array
func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := make([]string, 0, len(strSlice))
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// ReplaceAtIndex replaces a rune at index in string
func ReplaceAtIndex(in string, r rune, i int) string {
	if i < 0 || i >= len(in) {
		CheckError(fmt.Errorf("index out of range"))
	}
	out := []rune(in)
	if i >= 0 && i < len(out) {
		out[i] = r
		// In instances where i is out of bounds go to the end
	} else if i >= 0 && i == len(out) {
		out[len(out)-1] = r
	}
	return string(out)
}

// ReplaceWord replaces a mask within an input string with a value
func ReplaceWord(stringword, mask string, value string, replacements []string) string {
	tokenmask := MakeMask(value, replacements)
	tokenmask = models.ValidateMask(tokenmask)

	if strings.Contains(mask, tokenmask) {
		newword := strings.Replace(mask, tokenmask, value, -1)
		newword = strings.NewReplacer("?u", "?", "?l", "?", "?b", "?", "?d", "?", "?s", "?").Replace(newword)

		for i := 0; i < len(stringword); {
			r, size := utf8.DecodeRuneInString(stringword[i:])
			if i < len(newword) {
				if newword[i] == '?' {
					newword = ReplaceAtIndex(newword, r, i)
				}
			}
			i += size
		}

		// NOTE: This introduces a known bug
		// If the first string contains "?" and a multibyte character the
		// output is malformed
		if !strings.Contains(stringword, "?") {
			newword = strings.ReplaceAll(newword, "?", "")
		}

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
		os.Exit(0)
	}
}

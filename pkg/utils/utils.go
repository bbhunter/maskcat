// Package utils contains functions for the main maskcat program
package utils

import (
	"fmt"
	"os"
	"strings"
)

// MakeMask performs substitution to make HC masks
func MakeMask(str string) string {
	replacer := strings.NewReplacer(
		"a", "?l",
		"b", "?l",
		"c", "?l",
		"d", "?l",
		"e", "?l",
		"f", "?l",
		"g", "?l",
		"h", "?l",
		"i", "?l",
		"j", "?l",
		"k", "?l",
		"l", "?l",
		"m", "?l",
		"n", "?l",
		"o", "?l",
		"p", "?l",
		"q", "?l",
		"r", "?l",
		"s", "?l",
		"t", "?l",
		"u", "?l",
		"v", "?l",
		"w", "?l",
		"x", "?l",
		"y", "?l",
		"z", "?l",
		"A", "?u",
		"B", "?u",
		"C", "?u",
		"D", "?u",
		"E", "?u",
		"F", "?u",
		"G", "?u",
		"H", "?u",
		"I", "?u",
		"J", "?u",
		"K", "?u",
		"L", "?u",
		"M", "?u",
		"N", "?u",
		"O", "?u",
		"P", "?u",
		"Q", "?u",
		"R", "?u",
		"S", "?u",
		"T", "?u",
		"U", "?u",
		"V", "?u",
		"W", "?u",
		"X", "?u",
		"Y", "?u",
		"Z", "?u",
		"0", "?d",
		"1", "?d",
		"2", "?d",
		"3", "?d",
		"4", "?d",
		"5", "?d",
		"6", "?d",
		"7", "?d",
		"8", "?d",
		"9", "?d",
		" ", "?s",
		"!", "?s",
		"\"", "?s",
		"#", "?s",
		"$", "?s",
		"%", "?s",
		"&", "?s",
		"\\", "?s",
		"(", "?s",
		")", "?s",
		"*", "?s",
		"+", "?s",
		",", "?s",
		"-", "?s",
		".", "?s",
		"'", "?s",
		"/", "?s",
		":", "?s",
		";", "?s",
		"<", "?s",
		"=", "?s",
		">", "?s",
		"?", "?s",
		"@", "?s",
		"[", "?s",
		"]", "?s",
		"^", "?s",
		"_", "?s",
		"`", "?s",
		"{", "?s",
		"|", "?s",
		"}", "?s",
		"~", "?s",
	)
	return replacer.Replace(str)
}

// MakeToken replaces all non-alpha characters to generate tokens
func MakeToken(str string) string {
	replacer := strings.NewReplacer(
		"0", "",
		"1", "",
		"2", "",
		"3", "",
		"4", "",
		"5", "",
		"6", "",
		"7", "",
		"8", "",
		"9", "",
		"!", "",
		"\"", "",
		"#", "",
		"$", "",
		"%", "",
		"&", "",
		"\\", "",
		"(", "",
		")", "",
		"*", "",
		"+", "",
		",", "",
		"-", "",
		".", "",
		"'", "",
		"/", "",
		":", "",
		";", "",
		"<", "",
		"=", "",
		">", "",
		"?", "",
		"@", "",
		"[", "",
		"]", "",
		"^", "",
		"_", "",
		"`", "",
		"{", "",
		"|", "",
		"}", "",
		"~", "",
	)
	return replacer.Replace(str)
}

// MakePartialMask creates a
func MakePartialMask(str string, chars string) string {
	lowerReplacer := strings.NewReplacer(
		"a", "?l",
		"b", "?l",
		"c", "?l",
		"d", "?l",
		"e", "?l",
		"f", "?l",
		"g", "?l",
		"h", "?l",
		"i", "?l",
		"j", "?l",
		"k", "?l",
		"l", "?l",
		"m", "?l",
		"n", "?l",
		"o", "?l",
		"p", "?l",
		"q", "?l",
		"r", "?l",
		"s", "?l",
		"t", "?l",
		"u", "?l",
		"v", "?l",
		"w", "?l",
		"x", "?l",
		"y", "?l",
		"z", "?l")

	upperReplacer := strings.NewReplacer(
		"A", "?u",
		"B", "?u",
		"C", "?u",
		"D", "?u",
		"E", "?u",
		"F", "?u",
		"G", "?u",
		"H", "?u",
		"I", "?u",
		"J", "?u",
		"K", "?u",
		"L", "?u",
		"M", "?u",
		"N", "?u",
		"O", "?u",
		"P", "?u",
		"Q", "?u",
		"R", "?u",
		"S", "?u",
		"T", "?u",
		"U", "?u",
		"V", "?u",
		"W", "?u",
		"X", "?u",
		"Y", "?u",
		"Z", "?u")

	digitReplacer := strings.NewReplacer(
		"0", "?d",
		"1", "?d",
		"2", "?d",
		"3", "?d",
		"4", "?d",
		"5", "?d",
		"6", "?d",
		"7", "?d",
		"8", "?d",
		"9", "?d")

	specialReplacer := strings.NewReplacer(
		" ", "?s",
		"!", "?s",
		"\"", "?s",
		"#", "?s",
		"$", "?s",
		"%", "?s",
		"&", "?s",
		"\\", "?s",
		"(", "?s",
		")", "?s",
		"*", "?s",
		"+", "?s",
		",", "?s",
		"-", "?s",
		".", "?s",
		"'", "?s",
		"/", "?s",
		":", "?s",
		";", "?s",
		"<", "?s",
		"=", "?s",
		">", "?s",
		"?", "?s",
		"@", "?s",
		"[", "?s",
		"]", "?s",
		"^", "?s",
		"_", "?s",
		"`", "?s",
		"{", "?s",
		"|", "?s",
		"}", "?s",
		"~", "?s")

	if strings.Contains(chars, "u") {
		str = upperReplacer.Replace(str)
	}

	if strings.Contains(chars, "l") {
		str = lowerReplacer.Replace(str)
	}

	if strings.Contains(chars, "d") {
		str = digitReplacer.Replace(str)
	}

	if strings.Contains(chars, "s") {
		str = specialReplacer.Replace(str)
	}

	return str
}

// TestComplexity tests the complexity of an input string
func TestComplexity(str string) int {
	c := 0
	if strings.Contains(str, "?u") {
		c++
	}
	if strings.Contains(str, "?l") {
		c++
	}
	if strings.Contains(str, "?d") {
		c++
	}
	if strings.Contains(str, "?s") {
		c++
	}
	return c
}

// TestEntropy calculates mask entropy
func TestEntropy(str string) int {
	c := 0
	if strings.Contains(str, "?u") {
		c += strings.Count(str, "?u") * 26
	}
	if strings.Contains(str, "?l") {
		c += strings.Count(str, "?l") * 26
	}
	if strings.Contains(str, "?d") {
		c += strings.Count(str, "?d") * 10
	}
	if strings.Contains(str, "?s") {
		c += strings.Count(str, "?s") * 33
	}
	return c
}

// ChunkString splits string into chunks
func ChunkString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

// RemoveDuplicateStr removes duplicate strings from array
func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
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
	out := []rune(in)
	if len(out) == i {
		i := i - 1
		out[i] = r
	} else if i <= len(out) {
		out[i] = r
	}
	return string(out)
}

// ReplaceWord replaces a mask within an input string with a value
func ReplaceWord(stringword, mask string, value string) string {
	tokenmask := MakeMask(value)
	if strings.Contains(mask, tokenmask) {

		// format mask token chars into sub chars
		newword := strings.Replace(mask, tokenmask, value, -1)
		newword = strings.Replace(newword, "?u", "?", -1)
		newword = strings.Replace(newword, "?l", "?", -1)
		newword = strings.Replace(newword, "?d", "?", -1)
		newword = strings.Replace(newword, "?s", "?", -1)

		// loop over the string and finish the sub
		for i, c := range stringword {
			for x, y := range newword {
				if string(y) == "?" && x == i {
					newword = ReplaceAtIndex(newword, rune(c), i)
					break
				}
			}

		}
		if strings.Contains(newword, value) {
			if newword != value {
				return newword
			}
		}
	}
	return ""
}

// CheckError is a general error handler
func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(0)
	}
}

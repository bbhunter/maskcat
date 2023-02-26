// Package utils contains functions for the main maskcat program
package utils

import (
	"fmt"
	"os"
	"strings"
)

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
func ReplaceWord(stringword, mask string, value string, replacer *strings.Replacer) string {
	tokenmask := replacer.Replace(value)
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

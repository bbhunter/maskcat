// Package models holds all of the data structures and validations
package models

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// IsMask tests a string to see if it is a valid mask
func IsMask(mask string) bool {
	var IsMask = regexp.MustCompile(`^[uldsb?]+$`).MatchString
	if IsMask(mask) == false {
		return false
	}
	return true
}

// IsMaskChars tests a string to see if it only contains valid mask characters
func IsMaskChars(ch string) bool {
	var IsMaskChars = regexp.MustCompile(`^[uldsb]+$`).MatchString
	if IsMaskChars(ch) == false {
		return false
	}
	return true
}

// IsInt tests a string to see if it only contains numerical characters
func IsInt(str string) bool {
	var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString
	if IsInt(str) == false {
		return false
	}
	return true
}

// IsAlpha tests a string to see if it only contains alpha characters
func IsAlpha(str string) bool {
	var IsToken = regexp.MustCompile(`^[a-zA-Z ]*$`).MatchString
	if IsToken(str) == false {
		return false
	}
	return true
}

// CheckASCIIString checks to see if a string only contains ascii characters
func CheckASCIIString(str string) bool {
	if utf8.RuneCountInString(str) != len(str) {
		return false
	}
	return true
}

// ValidateMask checks if a string is a valid mask
// If not attempts to convert multibyte chars to a valid mask
func ValidateMask(mask string) string {
	if IsMask(mask) == false {
		if !CheckASCIIString(mask) {
			mask = ConvertMultiByteString(mask)
		}
	}
	return mask
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

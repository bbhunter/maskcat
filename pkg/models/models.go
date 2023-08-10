// Package models holds all of the data structures and validations
//
// The package structure is broken into two components:
//
// models.go which contains the primary logic
// models_test.go which contains unit tests
package models

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// IsHashMask tests a string to see if it contains only mask characters
//
// Args:
//
//	mask (string): The input string
//
// Returns:
//
//	(bool): If the string is a valid mask
func IsHashMask(mask string) bool {
	var IsMask = regexp.MustCompile(`^[uldsb?]+$`).MatchString
	if IsMask(mask) == false {
		return false
	}
	return true
}

// IsStringInt tests a string to see if it only contains numerical characters
//
// Args:
//
//	str (string): The input string
//
// Returns:
//
//	(bool) : If the string is a valid integer
func IsStringInt(str string) bool {
	var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString
	if IsInt(str) == false {
		return false
	}
	return true
}

// IsStringAlpha tests a string to see if it only contains alpha characters
//
// Args:
//
//	str (string): The input string
//
// Returns:
//
//	(bool) : If the string is a valid alpha only
func IsStringAlpha(str string) bool {
	var IsToken = regexp.MustCompile(`^[a-zA-Z ]*$`).MatchString
	if IsToken(str) == false {
		return false
	}
	return true
}

// IsStringASCII checks to see if a string only contains ASCII characters
//
// Args:
//
//	str (string): The input string
//
// Returns:
//
//	(bool) : If the string is a valid ASCII only
func IsStringASCII(str string) bool {
	if utf8.RuneCountInString(str) != len(str) {
		return false
	}
	return true
}

// EnsureValidMask checks if a string is a valid mask and transforms it if not
//
// # Used by the application to convert multibyte chars to a valid mask
//
// Args:
//
//	mask (string): Input string as a mask
//
// Return:
//
//	mask (string): Valid mask string
func EnsureValidMask(mask string) string {
	if IsHashMask(mask) == false {
		if !IsStringASCII(mask) {
			mask = ConvertMultiByteString(mask)
		}
	}
	return mask
}

// ConvertMultiByteString converts non-ascii characters to a valid format
//
// Args:
//
//	str (string): Input string
//
// Returns:
//
//	returnStr (string): Converted string
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

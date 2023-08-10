package models

import "testing"

func TestIsHashMask(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"uldsb?", true},
		{"", false},
		{"abc", false},
		{"uldsb", true},
	}

	for _, test := range tests {
		result := IsHashMask(test.input)
		if result != test.expected {
			t.Errorf("IsMask(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestIsStringInt(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},
		{"", false},
		{"abc", false},
		{"12a34", false},
	}

	for _, test := range tests {
		result := IsStringInt(test.input)
		if result != test.expected {
			t.Errorf("IsInt(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestIsStringAlpha(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abcABC", true},
		{"", true},
		{"abc123", false},
		{"abc ABC", true},
	}

	for _, test := range tests {
		result := IsStringAlpha(test.input)
		if result != test.expected {
			t.Errorf("IsAlpha(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestIsStringASCII(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello, 世界", false},
		{"Hello, World!", true},
		{"", true},
		{"世界", false},
	}

	for _, test := range tests {
		result := IsStringASCII(test.input)
		if result != test.expected {
			t.Errorf("CheckASCIIString(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestEnsureValidMask(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"uldsb?", "uldsb?"},
		{"", ""},
		{"abc", "abc"},
		{"世界", "?b?b?b?b?b?b"},
	}

	for _, test := range tests {
		result := EnsureValidMask(test.input)
		if result != test.expected {
			t.Errorf("ValidateMask(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestConvertMultiByteString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, 世界", "Hello, ?b?b?b?b?b?b"},
		{"", ""},
		{"abc", "abc"},
		{"世", "?b?b?b"},
	}

	for _, test := range tests {
		result := ConvertMultiByteString(test.input)
		if result != test.expected {
			t.Errorf("ConvertMultiByteString(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

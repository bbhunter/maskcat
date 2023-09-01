package utils

import (
	"reflect"
	"testing"
)

func TestConstructReplacements(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want []string
	}{
		{
			name: "Test lower case",
			str:  "l",
			want: []string{"a", "?l", "b", "?l", "c", "?l", "d", "?l", "e", "?l", "f", "?l", "g", "?l", "h", "?l", "i", "?l", "j", "?l", "k", "?l", "l", "?l", "m", "?l", "n", "?l", "o", "?l", "p", "?l", "q", "?l", "r", "?l", "s", "?l", "t", "?l", "u", "?l", "v", "?l", "w", "?l", "x", "?l", "y", "?l", "z", "?l"},
		},
		{
			name: "Test upper case",
			str:  "u",
			want: []string{"A", "?u", "B", "?u", "C", "?u", "D", "?u", "E", "?u", "F", "?u", "G", "?u", "H", "?u", "I", "?u", "J", "?u", "K", "?u", "L", "?u", "M", "?u", "N", "?u", "O", "?u", "P", "?u", "Q", "?u", "R", "?u", "S", "?u", "T", "?u", "U", "?u", "V", "?u", "W", "?u", "X", "?u", "Y", "?u", "Z", "?u"},
		},
		{
			name: "Test digits",
			str:  "d",
			want: []string{"0", "?d", "1", "?d", "2", "?d", "3", "?d", "4", "?d", "5", "?d", "6", "?d", "7", "?d", "8", "?d", "9", "?d"},
		},
		{
			name: "Test special characters",
			str:  "s",
			want: []string{" ", "?s", "!", "?s", "\"", "?s", "#", "?s", "$", "?s", "%", "?s", "&", "?s", "\\", "?s", "(", "?s", ")", "?s", "*", "?s", "+", "?s", ",", "?s", "-", "?s", ".", "?s", "/", "?s", ":", "?s", ";", "?s", "<", "?s", "=", "?s", ">", "?s", "?", "?s", "@", "?s", "[", "?s", "\\", "?s", "]", "?s", "^", "?s", "_", "?s", "`", "?s", "{", "?s", "|", "?s", "}", "?s", "~", "?s", "'", "?s"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstructReplacements(tt.str)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructReplacements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeMask(t *testing.T) {
	str := "Hello, World1!"
	replacements := ConstructReplacements("ulds")
	want := "?u?l?l?l?l?s?s?u?l?l?l?l?d?s"
	got := MakeMask(str, replacements)
	if got != want {
		t.Errorf("MakeMask(%q, %q) = %q; want %q", str, replacements, got, want)
	}
}

func TestMakeToken(t *testing.T) {
	str := "ThisApple123OfMine"
	want := []string{"This", "Apple", "123", "Of", "Mine", "ThisAppleOfMine"}
	got := MakeToken(str)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MakeToken(%q) = %q; want %q", str, got, want)
	}
}

func TestTestComplexity(t *testing.T) {
	str := "?u?l?l?l?l?s?s?u?l?l?l?l?d?s"
	want := 4
	got := TestComplexity(str)
	if got != want {
		t.Errorf("TestComplexity(%q) = %d; want %d", str, got, want)
	}
}

func TestTestEntropy(t *testing.T) {
	str := "?u?l?l?l?l?s?s?u?l?l?l?l?d?s"
	want := 369
	got := TestEntropy(str)
	if got != want {
		t.Errorf("TestEntropy(%q) = %d; want %d", str, got, want)
	}
}

func TestChunkString(t *testing.T) {
	s := "Hello, World!"
	chunkSize := 5
	want := []string{"Hello", ", Wor", "ld!"}
	got := ChunkString(s, chunkSize)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ChunkString(%q, %d) = %q; want %q", s, chunkSize, got, want)
	}
}

func TestReplaceWord(t *testing.T) {
	stringword := "Bello Jello Mello"
	mask := "?u?l?l?l?l?s?u?l?l?l?l?s?u?l?l?l?l"
	value := "Hello"
	replacements := ConstructReplacements("ulds")
	want := "Hello Jello Mello"
	got := ReplaceWord(stringword, mask, value, replacements, 1, 0)
	if got != want {
		t.Errorf("ReplaceWord(%q, %q, %q) = %q; want %q", stringword, mask, value, got, want)
	}
}

func TestRemoveMaskChars(t *testing.T) {
	str := "?u?l?d?s"
	want := ""
	got := RemoveMaskCharacters(str)
	if got != want {
		t.Errorf("RemoveMaskChars(%q) = %q; want %q", str, got, want)
	}
}

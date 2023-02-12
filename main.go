// Package that contains all relevant code for maskcat
// Note that ?b is not supported at this time
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
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

	stdinscanner := bufio.NewScanner(os.Stdin)

	if len(os.Args) > 1 {

		if len(os.Args) <= 2 {
			fmt.Println("OPTIONS: match sub mutate")
			fmt.Println("EXAMPLE: stdin | maskcat match masks.lst")
			fmt.Println("EXAMPLE: stdin | maskcat sub tokens.lst")
			fmt.Println("EXAMPLE: stdin | maskcat mutate <MAX-TOKEN-LEN>")
			os.Exit(0)
		}

		if os.Args[1] == "match" {
			infile := os.Args[2]
			buf, err := os.Open(infile)
			CheckError(err)

			defer func() {
				if err = buf.Close(); err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}()

			filescanner := bufio.NewScanner(buf)
			var masks []string

			for filescanner.Scan() {
				var IsMask = regexp.MustCompile(`^[ulds?]+$`).MatchString
				if IsMask(filescanner.Text()) == false {
					fmt.Println("[SKIP] Input mask contains non-mask characters: ", filescanner.Text())
					continue
				}
				masks = append(masks, filescanner.Text())
			}

			for stdinscanner.Scan() {
				mask := replacer.Replace(stdinscanner.Text())

				for _, value := range masks {

					if mask == value {
						fmt.Println(stdinscanner.Text())
						break
					}

					if err := stdinscanner.Err(); err != nil {
						fmt.Fprintln(os.Stderr, "reading standard input:", err)
					}
				}
			}

		} else if os.Args[1] == "sub" {
			infile := os.Args[2]
			buf, err := os.Open(infile)
			CheckError(err)

			defer func() {
				if err = buf.Close(); err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}()

			filescanner := bufio.NewScanner(buf)
			var tokens []string

			for filescanner.Scan() {
				if filescanner.Text() != "" {
					tokens = append(tokens, filescanner.Text())
				}

				if err := filescanner.Err(); err != nil {
					fmt.Fprintln(os.Stderr, "reading standard input:", err)
				}
			}

			for stdinscanner.Scan() {
				stringword := stdinscanner.Text()
				mask := replacer.Replace(stdinscanner.Text())

				for _, value := range tokens {
					newWord := ReplaceWord(stringword, mask, value, replacer)

					if newWord != "" {
						fmt.Println(newWord)
					}
				}
			}

		} else if os.Args[1] == "mutate" {
			var tokens []string
			var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString

			if IsInt(os.Args[2]) == false {
				CheckError(errors.New("ERROR: Invalid Chunk Size"))
			}

			for stdinscanner.Scan() {
				chunksInt, err := strconv.Atoi(os.Args[2])
				CheckError(err)
				chunks := ChunkString(stdinscanner.Text(), chunksInt)
				var achunks []string
				for _, ch := range chunks {
					if len(ch) == chunksInt {
						achunks = append(achunks, ch)
					}
				}
				tokens = append(tokens, achunks...)
				tokens = RemoveDuplicateStr(tokens)

				stringword := stdinscanner.Text()
				mask := replacer.Replace(stdinscanner.Text())

				for _, value := range tokens {
					newWord := ReplaceWord(stringword, mask, value, replacer)

					if newWord != "" {
						fmt.Println(newWord)
					}
				}
			}
		}

	} else {
		// else make masks
		for stdinscanner.Scan() {
			mask := replacer.Replace(stdinscanner.Text())
			var IsMask = regexp.MustCompile(`^[ulds?]+$`).MatchString
			if IsMask(mask) == false {
				continue
			}
			fmt.Printf("%s:%d:%d:%d\n", mask, len(stdinscanner.Text()), TestComplexity(mask), TestEntropy(mask))
		}

		if err := stdinscanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
	for scanner.Scan() {
		mask := replacer.Replace(scanner.Text())
		fmt.Printf("%s:%d:%d\n", mask, len(scanner.Text()), TestComplexity(mask))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

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

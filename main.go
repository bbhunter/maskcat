package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
			fmt.Println("OPTIONS: match sub")
			fmt.Println("EXAMPLE: stdin | maskcat match masks.lst")
			fmt.Println("EXAMPLE: stdin | maskcat sub tokens.lst")
			os.Exit(0)
		}
		infile := os.Args[2]

		buf, err := os.Open(infile)
		checkError(err)

		defer func() {
			if err = buf.Close(); err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}()

		filescanner := bufio.NewScanner(buf)

		if os.Args[1] == "match" || os.Args[1] == "m" {
			var masks []string

			for filescanner.Scan() {
				var IsMask = regexp.MustCompile(`^[ulds?]+$`).MatchString
				if IsMask(filescanner.Text()) == false {
					//fmt.Println("Input mask contains non-mask characters. Passing.")
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
		} else if os.Args[1] == "sub" || os.Args[1] == "s" {
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
				mask := replacer.Replace(stdinscanner.Text())

				for _, value := range tokens {
					tokenmask := replacer.Replace(value)
					if strings.Contains(mask, tokenmask) {

						// format mask token chars into sub chars
						newword := strings.Replace(mask, tokenmask, value, -1)
						newword = strings.Replace(newword, "?u", "?", -1)
						newword = strings.Replace(newword, "?l", "?", -1)
						newword = strings.Replace(newword, "?d", "?", -1)
						newword = strings.Replace(newword, "?s", "?", -1)

						// loop over the string and finish the sub
						for i, c := range stdinscanner.Text() {
							for x, y := range newword {
								if string(y) == "?" && x == i {
									newword = replaceAtIndex(newword, rune(c), i)
									break
								}
							}

						}
						if strings.Contains(newword, value) {
							if newword != value {
								fmt.Println(newword)
							}
						}
					}

					if err := stdinscanner.Err(); err != nil {
						fmt.Fprintln(os.Stderr, "reading standard input:", err)
					}
				}

			}

		}

	} else {

		for stdinscanner.Scan() {
			mask := replacer.Replace(stdinscanner.Text())
			// check to see if the mask contains invalid characters. if so pass
			var IsMask = regexp.MustCompile(`^[ulds?]+$`).MatchString
			if IsMask(mask) == false {
				continue
			}
			fmt.Printf("%s:%d:%d\n", mask, len(stdinscanner.Text()), TestComplexity(mask))
		}

		if err := stdinscanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

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

// Replace rune at index in string
func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

// Error checking
func checkError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(0)
	}
}

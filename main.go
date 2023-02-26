// Package that contains the primary logic for maskcat and the CLI
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

	"github.com/jakewnuk/maskcat/pkg/utils"
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
			utils.CheckError(err)

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
			utils.CheckError(err)

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
					newWord := utils.ReplaceWord(stringword, mask, value, replacer)

					if newWord != "" {
						fmt.Println(newWord)
					}
				}
			}

		} else if os.Args[1] == "mutate" {
			var tokens []string
			var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString

			if IsInt(os.Args[2]) == false {
				utils.CheckError(errors.New("ERROR: Invalid Chunk Size"))
			}

			for stdinscanner.Scan() {
				chunksInt, err := strconv.Atoi(os.Args[2])
				utils.CheckError(err)
				chunks := utils.ChunkString(stdinscanner.Text(), chunksInt)
				var achunks []string
				for _, ch := range chunks {
					if len(ch) == chunksInt {
						achunks = append(achunks, ch)
					}
				}
				tokens = append(tokens, achunks...)
				tokens = utils.RemoveDuplicateStr(tokens)

				stringword := stdinscanner.Text()
				mask := replacer.Replace(stdinscanner.Text())

				for _, value := range tokens {
					newWord := utils.ReplaceWord(stringword, mask, value, replacer)

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
			fmt.Printf("%s:%d:%d:%d\n", mask, len(stdinscanner.Text()), utils.TestComplexity(mask), utils.TestEntropy(mask))
		}

		if err := stdinscanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
}

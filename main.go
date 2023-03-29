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

	"github.com/jakewnuk/maskcat/pkg/utils"
)

func main() {

	stdIn := bufio.NewScanner(os.Stdin)

	if len(os.Args) > 1 {

		if len(os.Args) <= 2 {
			fmt.Println("OPTIONS: match sub mutate tokens partial")
			fmt.Println("EXAMPLE: stdin | maskcat match masks.lst")
			fmt.Println("EXAMPLE: stdin | maskcat sub tokens.lst")
			fmt.Println("EXAMPLE: stdin | maskcat mutate <MAX-TOKEN-LEN>")
			fmt.Println("EXAMPLE: stdin | maskcat tokens <MAX-LEN> (values over 99 allow all)")
			fmt.Println("EXAMPLE: stdin | maskcat partial <MASK-CHARS>")
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

			for stdIn.Scan() {
				mask := utils.MakeMask(stdIn.Text())

				for _, value := range masks {

					if mask == value {
						fmt.Println(stdIn.Text())
						break
					}

					if err := stdIn.Err(); err != nil {
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

			for stdIn.Scan() {
				stringword := stdIn.Text()
				mask := utils.MakeMask(stdIn.Text())

				for _, value := range tokens {
					newWord := utils.ReplaceWord(stringword, mask, value)

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

			for stdIn.Scan() {
				chunksInt, err := strconv.Atoi(os.Args[2])
				utils.CheckError(err)
				chunks := utils.ChunkString(stdIn.Text(), chunksInt)
				var achunks []string
				for _, ch := range chunks {
					if len(ch) == chunksInt {
						achunks = append(achunks, ch)
					}
				}
				tokens = append(tokens, achunks...)
				tokens = utils.RemoveDuplicateStr(tokens)

				stringword := stdIn.Text()
				mask := utils.MakeMask(stdIn.Text())

				for _, value := range tokens {
					newWord := utils.ReplaceWord(stringword, mask, value)

					if newWord != "" {
						fmt.Println(newWord)
					}
				}
			}
		} else if os.Args[1] == "tokens" {
			var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString

			if IsInt(os.Args[2]) == false {
				utils.CheckError(errors.New("ERROR: Invalid String Size"))
			}

			for stdIn.Scan() {
				token := utils.MakeToken(stdIn.Text())
				var IsToken = regexp.MustCompile(`^[a-zA-Z ]*$`).MatchString
				if IsToken(token) == false {
					continue
				}

				length, err := strconv.Atoi(os.Args[2])
				utils.CheckError(err)

				// VALUES OVER 99 LET ALL THROUGH
				if len(token) != length && length < 98 {
					continue
				}

				fmt.Printf("%s\n", token)
			}
		} else if os.Args[1] == "partial" {
			var IsMaskChars = regexp.MustCompile(`^[ulds]+$`).MatchString

			if IsMaskChars(os.Args[2]) == false {
				utils.CheckError(errors.New("ERROR: Can only contain 'u','d','l', and 's'"))
			}

			for stdIn.Scan() {
				partial := utils.MakePartialMask(stdIn.Text(), os.Args[2])
				fmt.Printf("%s\n", partial)
			}

		}

	} else {
		// else make masks
		for stdIn.Scan() {
			mask := utils.MakeMask(stdIn.Text())
			var IsMask = regexp.MustCompile(`^[ulds?]+$`).MatchString
			if IsMask(mask) == false {
				continue
			}
			fmt.Printf("%s:%d:%d:%d\n", mask, len(stdIn.Text()), utils.TestComplexity(mask), utils.TestEntropy(mask))
		}

		if err := stdIn.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
}

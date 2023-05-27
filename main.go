// Package that contains the primary logic for maskcat and the CLI
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
	if len(os.Args) <= 2 {
		if len(os.Args) == 1 {
			os.Args = append(os.Args, "default")
		} else {
			printUsage()
			os.Exit(0)
		}
	}

	stdIn := bufio.NewScanner(os.Stdin)

	switch os.Args[1] {
	case "match":
		matchMasks(stdIn, os.Args[2])
	case "sub":
		subMasks(stdIn, os.Args[2])
	case "mutate":
		mutateMasks(stdIn, os.Args[2])
	case "tokens":
		generateTokens(stdIn, os.Args[2])
	case "partial":
		generatePartialMasks(stdIn, os.Args[2])
	case "remove":
		generatePartialRemoveMasks(stdIn, os.Args[2])
	default:
		generateMasks(stdIn)
	}
}

// printUsage prints usage information for the program
func printUsage() {
	fmt.Println("OPTIONS: match sub mutate tokens partial remove")
	fmt.Println("EXAMPLE: stdin | maskcat match <MASK-FILE>")
	fmt.Println("EXAMPLE: stdin | maskcat sub <TOKENS-FILE>")
	fmt.Println("EXAMPLE: stdin | maskcat mutate <MAX-TOKEN-LEN>")
	fmt.Println("EXAMPLE: stdin | maskcat tokens <TOKEN-LEN> (99+ returns all)")
	fmt.Println("EXAMPLE: stdin | maskcat partial <MASK-CHARS>")
	fmt.Println("EXAMPLE: stdin | maskcat remove <MASK-CHARS>")
}

// matchMasks reads masks from a file and prints any input strings that match one of the masks
func matchMasks(stdIn *bufio.Scanner, infile string) {
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
	args := utils.ConstructReplacements("ulds")

	for filescanner.Scan() {
		var IsMask = regexp.MustCompile(`^[ulds?]+$`).MatchString
		if IsMask(filescanner.Text()) == false {
			fmt.Println("[SKIP] Input mask contains non-mask characters: ", filescanner.Text())
			continue
		}
		masks = append(masks, filescanner.Text())
	}

	for stdIn.Scan() {
		mask := utils.MakeMask(stdIn.Text(), args)

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
}

// subMasks reads tokens from a file and replaces mask characters in the input strings with the tokens
func subMasks(stdIn *bufio.Scanner, infile string) {
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
	args := utils.ConstructReplacements("ulds")

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
		mask := utils.MakeMask(stdIn.Text(), args)

		for _, value := range tokens {
			newWord := utils.ReplaceWord(stringword, mask, value, args)

			if newWord != "" {
				fmt.Println(newWord)
			}
		}
	}
}

// mutateMasks splits the input strings into chunks and replaces mask characters with the chunks
func mutateMasks(stdIn *bufio.Scanner, chunkSizeStr string) {
	var tokens []string
	var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString
	args := utils.ConstructReplacements("ulds")

	if IsInt(chunkSizeStr) == false {
		utils.CheckError(errors.New("Invalid Chunk Size"))
	}

	for stdIn.Scan() {
		chunksInt, err := strconv.Atoi(chunkSizeStr)
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
		mask := utils.MakeMask(stdIn.Text(), args)

		for _, value := range tokens {
			newWord := utils.ReplaceWord(stringword, mask, value, args)

			if newWord != "" {
				fmt.Println(newWord)
			}
		}
	}
}

// generateTokens generates tokens from the input strings by removing all non-alpha characters
func generateTokens(stdIn *bufio.Scanner, lengthStr string) {
	var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString

	if IsInt(lengthStr) == false {
		utils.CheckError(errors.New("Invalid String Size"))
	}

	for stdIn.Scan() {
		token := utils.MakeToken(stdIn.Text())
		var IsToken = regexp.MustCompile(`^[a-zA-Z ]*$`).MatchString
		if IsToken(token) == false {
			continue
		}

		length, err := strconv.Atoi(lengthStr)
		utils.CheckError(err)

		// VALUES OVER 99 LET ALL THROUGH
		if len(token) != length && length < 98 {
			continue
		}

		fmt.Printf("%s\n", token)
	}
}

// generatePartialMasks generates partial masks from the input strings using the specified mask characters
func generatePartialMasks(stdIn *bufio.Scanner, maskChars string) {
	var IsMaskChars = regexp.MustCompile(`^[ulds]+$`).MatchString
	args := utils.ConstructReplacements(maskChars)

	if IsMaskChars(maskChars) == false {
		utils.CheckError(errors.New("Can only contain 'u','d','l', and 's'"))
	}

	for stdIn.Scan() {
		partial := utils.MakePartialMask(stdIn.Text(), args)
		fmt.Printf("%s\n", partial)
	}
}

// generatePartialRemoveMasks removes characters in masks from the input strings using the specified mask characters
func generatePartialRemoveMasks(stdIn *bufio.Scanner, maskChars string) {
	var IsMaskChars = regexp.MustCompile(`^[ulds]+$`).MatchString
	args := utils.ConstructReplacements(maskChars)

	if IsMaskChars(maskChars) == false {
		utils.CheckError(errors.New("Can only contain 'u','d','l', and 's'"))
	}

	for stdIn.Scan() {
		partial := utils.MakePartialMask(stdIn.Text(), args)
		remaining := utils.RemoveMaskChars(partial)
		fmt.Printf("%s\n", remaining)
	}
}

// generateMasks generates masks from the input strings and prints information about the masks
func generateMasks(stdIn *bufio.Scanner) {
	args := utils.ConstructReplacements("ulds")
	for stdIn.Scan() {
		mask := utils.MakeMask(stdIn.Text(), args)
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

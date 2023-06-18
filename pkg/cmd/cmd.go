// Package cmd that contains the primary CLI logic for maskcat
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jakewnuk/maskcat/pkg/models"
	"github.com/jakewnuk/maskcat/pkg/utils"
)

// MatchMasks reads masks from a file and prints any input strings that match one of the masks
func MatchMasks(stdIn *bufio.Scanner, infile string) {
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
		if models.IsMask(filescanner.Text()) == false {
			fmt.Println("[SKIP] Input mask contains non-mask characters: ", filescanner.Text())
			continue
		}
		masks = append(masks, filescanner.Text())
	}

	for stdIn.Scan() {
		mask := utils.MakeMask(stdIn.Text(), args)
		mask = models.ValidateMask(mask)

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

// SubMasks reads tokens from a file and replaces mask characters in the input strings with the tokens
func SubMasks(stdIn *bufio.Scanner, infile string) {
	buf, err := os.Open(infile)
	utils.CheckError(err)

	defer func() {
		if err = buf.Close(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	filescanner := bufio.NewScanner(buf)
	tokens := make(map[string]struct{})
	args := utils.ConstructReplacements("ulds")

	for filescanner.Scan() {
		if filescanner.Text() != "" {
			tokens[filescanner.Text()] = struct{}{}
		}

		if err := filescanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}

	for stdIn.Scan() {
		stringword := stdIn.Text()
		mask := utils.MakeMask(stdIn.Text(), args)
		mask = models.ValidateMask(mask)

		for value := range tokens {
			newWord := utils.ReplaceWord(stringword, mask, value, args)
			if newWord != "" {
				fmt.Println(newWord)
			}
		}
	}
}

// MutateMasks splits the input strings into chunks and replaces mask characters with the chunks
func MutateMasks(stdIn *bufio.Scanner, chunkSizeStr string) {
	tokens := make(map[string]struct{})
	args := utils.ConstructReplacements("ulds")

	if models.IsInt(chunkSizeStr) == false {
		utils.CheckError(errors.New("Invalid Chunk Size"))
	}

	for stdIn.Scan() {
		chunksInt, err := strconv.Atoi(chunkSizeStr)
		utils.CheckError(err)
		chunks := utils.ChunkString(stdIn.Text(), chunksInt)

		for _, ch := range chunks {
			if len(ch) == chunksInt {
				tokens[ch] = struct{}{}
			}
		}

		stringword := stdIn.Text()
		mask := utils.MakeMask(stdIn.Text(), args)
		mask = models.ValidateMask(mask)

		for value := range tokens {
			newWord := utils.ReplaceWord(stringword, mask, value, args)
			if newWord != "" {
				fmt.Println(newWord)
			}
		}
	}
}

// GenerateTokens generates tokens from the input strings by removing all non-alpha characters
func GenerateTokens(stdIn *bufio.Scanner, lengthStr string) {
	if models.IsInt(lengthStr) == false {
		utils.CheckError(errors.New("Invalid String Size"))
	}

	for stdIn.Scan() {
		token := utils.MakeToken(stdIn.Text())
		if models.IsAlpha(token) == false {
			continue
		}

		length, err := strconv.Atoi(lengthStr)
		utils.CheckError(err)

		// NOTE: VALUES OVER 99 LET ALL THROUGH
		if len(token) != length && length < 98 {
			continue
		}

		fmt.Printf("%s\n", token)
	}
}

// GeneratePartialMasks generates partial masks from the input strings using the specified mask characters
func GeneratePartialMasks(stdIn *bufio.Scanner, maskChars string) {
	args := utils.ConstructReplacements(maskChars)

	if models.IsMaskChars(maskChars) == false {
		utils.CheckError(errors.New("Can only contain 'u','d','l', 'b', and 's'"))
	}

	for stdIn.Scan() {
		partial := utils.MakePartialMask(stdIn.Text(), args)
		if strings.Contains(maskChars, "b") {
			partial = models.ConvertMultiByteString(partial)
		}
		fmt.Printf("%s\n", partial)
	}
}

// GeneratePartialRemoveMasks removes characters in masks from the input strings using the specified mask characters
func GeneratePartialRemoveMasks(stdIn *bufio.Scanner, maskChars string) {
	args := utils.ConstructReplacements(maskChars)

	if models.IsMaskChars(maskChars) == false {
		utils.CheckError(errors.New("Can only contain 'u','d','l', 'b', and 's'"))
	}

	for stdIn.Scan() {
		partial := utils.MakePartialMask(stdIn.Text(), args)
		if strings.Contains(maskChars, "b") {
			partial = models.ConvertMultiByteString(partial)
		}
		remaining := utils.RemoveMaskChars(partial)
		fmt.Printf("%s\n", remaining)
	}
}

// GenerateMasks generates masks from the input strings and prints information about the masks
func GenerateMasks(stdIn *bufio.Scanner) {
	args := utils.ConstructReplacements("ulds")
	for stdIn.Scan() {
		mask := utils.MakeMask(stdIn.Text(), args)
		mask = models.ValidateMask(mask)
		fmt.Printf("%s:%d:%d:%d\n", mask, len(stdIn.Text()), utils.TestComplexity(mask), utils.TestEntropy(mask))
	}

	if err := stdIn.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// Package cli contains logic for operating the cli tool
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jakewnuk/maskcat/pkg/models"
	"github.com/jakewnuk/maskcat/pkg/utils"
)

// MatchMasks reads masks from a file and prints any input strings that match one of the masks
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	infile (string): File path of input file to use
//
// Returns:
//
//	None
func MatchMasks(stdIn *bufio.Scanner, infile string, doMultiByte bool) {
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
	args := utils.ConstructReplacements("ulds")

	for filescanner.Scan() {
		if models.IsHashMask(filescanner.Text()) == false {
			fmt.Println("[SKIP] Input mask contains non-mask characters: ", filescanner.Text())
			continue
		}
		masks = append(masks, filescanner.Text())
	}

	results := make(chan string)
	var wg sync.WaitGroup

	for stdIn.Scan() {
		mask := utils.MakeMask(stdIn.Text(), args)

		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, value := range masks {

				if mask == value {
					results <- stdIn.Text()
					break
				}

				if err := stdIn.Err(); err != nil {
					CheckError(err)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

// SubMasks reads tokens from a file and replaces mask characters in the input strings with the tokens
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	infile (string): File path of input file to use
//
// Returns:
//
// None
func SubMasks(stdIn *bufio.Scanner, infile string, doMultiByte bool, doNumberOfReplacements int) {
	buf, err := os.Open(infile)
	CheckError(err)

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
			CheckError(err)
		}
	}

	results := make(chan string)
	var wg sync.WaitGroup

	for stdIn.Scan() {
		stringWord := stdIn.Text()
		mask := utils.MakeMask(stdIn.Text(), args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for value := range tokens {
				newWord := utils.ReplaceWord(stringWord, mask, value, args, doNumberOfReplacements)
				if newWord != "" {
					results <- newWord
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

// MutateMasks splits the input strings into chunks and replaces mask characters with the chunks
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	chunkSizeStr (string): Size of the chunks as a number
//
// Returns:
//
// None
func MutateMasks(stdIn *bufio.Scanner, chunkSizeStr string, doMultiByte bool, doNumberOfReplacements int) {
	var tokens sync.Map
	args := utils.ConstructReplacements("ulds")

	if models.IsStringInt(chunkSizeStr) == false {
		CheckError(errors.New("Invalid Chunk Size"))
	}

	var wg sync.WaitGroup

	for stdIn.Scan() {
		chunksInt, err := strconv.Atoi(chunkSizeStr)
		CheckError(err)
		chunks := utils.ChunkString(stdIn.Text(), chunksInt)

		for _, ch := range chunks {
			if len(ch) == chunksInt {
				tokens.Store(ch, struct{}{})
			}
		}

		stringWord := stdIn.Text()
		mask := utils.MakeMask(stdIn.Text(), args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			tokens.Range(func(key, value interface{}) bool {
				newWord := utils.ReplaceWord(stringWord, mask, key.(string), args, doNumberOfReplacements)
				if newWord != "" {
					fmt.Println(newWord)
				}
				return true
			})
		}()
	}
	wg.Wait()
}

// GenerateTokens generates tokens from the input strings by removing all non-alpha characters
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	lengthStr (string): Length of the tokens as a number
//
// Returns:
//
// None
func GenerateTokens(stdIn *bufio.Scanner, lengthStr string) {
	if models.IsStringInt(lengthStr) == false {
		CheckError(errors.New("Invalid String Size"))
	}

	for stdIn.Scan() {
		tokens := utils.MakeToken(stdIn.Text())
		for _, token := range tokens {
			if models.IsStringAlpha(token) == false {
				continue
			}

			length, err := strconv.Atoi(lengthStr)
			CheckError(err)

			// NOTE: VALUES OVER 99 LET ALL THROUGH
			if len(token) != length && length < 98 {
				continue
			}

			fmt.Printf("%s\n", token)
		}
	}
}

// GeneratePartialMasks generates partial masks from the input strings using the specified mask characters
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	maskChars (string): String of which character sets to replace (udlsb)
//
// Returns:
//
// None
func GeneratePartialMasks(stdIn *bufio.Scanner, maskChars string) {
	args := utils.ConstructReplacements(maskChars)

	if models.IsHashMask(maskChars) == false {
		CheckError(errors.New("Can only contain 'u','d','l', 'b', and 's'"))
	}

	for stdIn.Scan() {
		partial := utils.MakeMask(stdIn.Text(), args)
		if strings.Contains(maskChars, "b") {
			partial = models.ConvertMultiByteString(partial)
		}
		fmt.Printf("%s\n", partial)
	}
}

// GeneratePartialRemoveMasks removes characters in masks from the input strings using the specified mask characters
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	infile (string): File path of input file to use
//	maskChars (string): String of which character sets to replace (udlsb)
//
// Returns:
//
// None
func GeneratePartialRemoveMasks(stdIn *bufio.Scanner, maskChars string) {
	args := utils.ConstructReplacements(maskChars)

	if models.IsHashMask(maskChars) == false {
		CheckError(errors.New("Can only contain 'u','d','l', 'b', and 's'"))
	}

	for stdIn.Scan() {
		partial := utils.MakeMask(stdIn.Text(), args)
		if strings.Contains(maskChars, "b") {
			partial = models.ConvertMultiByteString(partial)
		}
		remaining := utils.RemoveMaskCharacters(partial)
		fmt.Printf("%s\n", remaining)
	}
}

// GenerateMasks generates masks from the input strings and prints information about the masks
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//
// Returns:
//
// None
func GenerateMasks(stdIn *bufio.Scanner, doMultiByte bool, verbose bool) {
	args := utils.ConstructReplacements("ulds")
	for stdIn.Scan() {
		mask := utils.MakeMask(stdIn.Text(), args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}
		if verbose {
			fmt.Printf("%s:%d:%d:%d\n", mask, len(stdIn.Text()), utils.TestComplexity(mask), utils.TestEntropy(mask))
		} else {
			fmt.Printf("%s\n", mask)
		}
	}

	if err := stdIn.Err(); err != nil {
		CheckError(err)
	}
}

// CheckIfArgExists checks an argument at a postion to see if it exists
//
// Args:
//
//	index (int): Index to check
//	args ([]string): Array of arguments
//
// Returns:
//
//	(bool): If the index exists
func CheckIfArgExists(index int, args []string) {
	exists := false
	if len(args) > index {
		exists = true
	}

	if !exists {
		CheckError(fmt.Errorf("Not enough arguments provided"))
	}
}

// CheckError is a general error handler
func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

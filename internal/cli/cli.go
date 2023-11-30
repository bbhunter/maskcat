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
//	doMultiByte (bool): If multibyte text should be processed
//	doDeHex (bool): If $HEX[...] text should be processed
//
// Returns:
//
//	None
func MatchMasks(stdIn *bufio.Scanner, infile string, doMultiByte bool, doDeHex bool) {
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
	stdText := ""

	for filescanner.Scan() {
		if models.IsHashMask(filescanner.Text()) == false {
			fmt.Println("[SKIP] Input mask contains non-mask characters: ", filescanner.Text())
			continue
		}
		masks = append(masks, filescanner.Text())
	}

	var wg sync.WaitGroup

	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				fmt.Println("error")
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		mask := utils.MakeMask(stdText, args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}

		wg.Add(1)
		go func(mask string, stdText string) {
			defer wg.Done()
			for _, value := range masks {

				if mask == value {
					fmt.Println(stdText)
					break
				}

				if err := stdIn.Err(); err != nil {
					CheckError(err)
				}
			}
		}(mask, stdText)
	}
	wg.Wait()
}

// SubMasks reads tokens from a file and replaces mask characters in the input strings with the tokens
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	infile (string): File path of input file to use
//	doMultiByte (bool): If multibyte text should be processed
//	doDeHex (bool): If $HEX[...] text should be processed
//	doNumberOfReplacements (int): Max number of times to replace per string
//	doFuzzAmount(int): Number of additional fuzz characters to add to replacer
//
// Returns:
//
// None
func SubMasks(stdIn *bufio.Scanner, infile string, doMultiByte bool, doDeHex bool, doNumberOfReplacements int, doFuzzAmount int) {
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

	var wg sync.WaitGroup

	for stdIn.Scan() {
		stringWord := ""
		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stringWord = ""
			}
			stringWord = plaintext
		} else {
			stringWord = stdIn.Text()
		}

		mask := utils.MakeMask(stringWord, args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}

		wg.Add(1)
		go func(stringWord string, mask string) {
			defer wg.Done()
			for value := range tokens {
				newWord := utils.ReplaceWordByMask(stringWord, mask, value, args, doNumberOfReplacements, doFuzzAmount)

				if newWord != "" {
					fmt.Println(newWord)
				}
			}
		}(stringWord, mask)
	}
	wg.Wait()
}

// MutateMasks splits the input strings into chunks and replaces mask characters with the chunks
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	chunkSizeStr (string): Size of the chunks as a number
//	doMultiByte (bool): If multibyte text should be processed
//	doDeHex (bool): If $HEX[...] text should be processed
//	doNumberOfReplacements (int): Max number of times to replace per string
//	doFuzzAmount(int): Number of additional fuzz characters to add to replacer
//
// Returns:
//
// None
func MutateMasks(stdIn *bufio.Scanner, chunkSizeStr string, doMultiByte bool, doDeHex bool, doNumberOfReplacements int, doFuzzAmount int) {
	var tokens sync.Map
	args := utils.ConstructReplacements("ulds")
	stdText := ""
	if models.IsStringInt(chunkSizeStr) == false {
		CheckError(errors.New("Invalid Chunk Size"))
	}

	var wg sync.WaitGroup

	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		ngrams := utils.MakeToken(stdText)
		chunksInt, err := strconv.Atoi(chunkSizeStr)
		CheckError(err)

		for _, token := range ngrams {
			if len(token) >= chunksInt {
				tokens.Store(token, struct{}{})
			}
		}

		stringWord := stdText
		mask := utils.MakeMask(stdText, args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}

		wg.Add(1)
		go func(stringWord string, mask string) {
			defer wg.Done()

			tokens.Range(func(key, value interface{}) bool {
				newWord := utils.ReplaceWordByMask(stringWord, mask, key.(string), args, doNumberOfReplacements, doFuzzAmount)
				if newWord != "" {
					fmt.Println(newWord)
				}
				return true
			})
		}(stringWord, mask)
	}
	wg.Wait()
}

// GenerateTokens generates tokens from the input strings by removing all non-alpha characters
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	lengthStr (string): Length of the tokens as a number
//	doDeHex (bool): If $HEX[...] text should be processed
//
// Returns:
//
// None
func GenerateTokens(stdIn *bufio.Scanner, lengthStr string, doDeHex bool) {
	stdText := ""
	if models.IsStringInt(lengthStr) == false {
		CheckError(errors.New("Invalid String Size"))
	}

	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		tokens := utils.MakeToken(stdText)
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
//	doDeHex (bool): If $HEX[...] text should be processed
//
// Returns:
//
// None
func GeneratePartialMasks(stdIn *bufio.Scanner, maskChars string, doDeHex bool) {
	args := utils.ConstructReplacements(maskChars)
	stdText := ""

	if models.IsHashMask(maskChars) == false {
		CheckError(errors.New("Can only contain 'u','d','l', 'b', and 's'"))
	}

	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		partial := utils.MakeMask(stdText, args)
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
//	doDeHex (bool): If $HEX[...] text should be processed
//
// Returns:
//
// None
func GeneratePartialRemoveMasks(stdIn *bufio.Scanner, maskChars string, doDeHex bool) {
	args := utils.ConstructReplacements(maskChars)
	stdText := ""

	if models.IsHashMask(maskChars) == false {
		CheckError(errors.New("Can only contain 'u','d','l', 'b', and 's'"))
	}

	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		partial := utils.MakeMask(stdText, args)
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
//	doMultiByte (bool): If multibyte text should be processed
//	doDeHex (bool): If $HEX[...] text should be processed
//	verbose (bool): If verbose stdText should be printed about masks
//
// Returns:
//
// None
func GenerateMasks(stdIn *bufio.Scanner, doMultiByte bool, doDeHex bool, verbose bool) {
	args := utils.ConstructReplacements("ulds")
	stdText := ""
	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		mask := utils.MakeMask(stdText, args)
		if doMultiByte {
			mask = models.EnsureValidMask(mask)
		}
		if verbose {
			fmt.Printf("%s:%d:%d:%d\n", mask, len(stdText), utils.TestComplexity(mask), utils.TestEntropy(mask))
		} else {
			fmt.Printf("%s\n", mask)
		}
	}

	if err := stdIn.Err(); err != nil {
		CheckError(err)
	}
}

// GenerateTokenRetainMasks creates masks while retaining tokens from a file
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	infile (string): File path of input file to use
//	doMultiByte (bool): If multibyte text should be processed
//	doDeHex (bool): If $HEX[...] text should be processed
//	doNumberOfReplacements (int): Max number of times to replace per string
//
// Returns:
//
// None
func GenerateTokenRetainMasks(stdIn *bufio.Scanner, infile string, doMultiByte bool, doDeHex bool, doNumberOfReplacements int) {
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

	var wg sync.WaitGroup

	for stdIn.Scan() {
		stringWord := ""
		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stringWord = ""
			}
			stringWord = plaintext
		} else {
			stringWord = stdIn.Text()
		}

		wg.Add(1)
		go func(stringWord string) {
			defer wg.Done()

			// Create the retain mask
			mask := utils.CreateRetainMask(stringWord, tokens, args, doMultiByte, doNumberOfReplacements)
			fmt.Println(mask)

		}(stringWord)
	}
	wg.Wait()
}

// GenerateSpliceMutation performs mutation mode on retain masks
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	infile (string): File path of input file to use
//	doMultiByte (bool): If multibyte text should be processed
//	doDeHex (bool): If $HEX[...] text should be processed
//	doNumberOfReplacements (int): Max number of times to replace per string
//	doFuzzAmount(int): Number of additional fuzz characters to add to replacer
//
// Returns:
//
// None
func GenerateSpliceMutation(stdIn *bufio.Scanner, infile string, doMultiByte bool, doDeHex bool, doNumberOfReplacements int, doFuzzAmount int) {

	// Read the retain infile
	buf, err := os.Open(infile)
	CheckError(err)

	defer func() {
		if err = buf.Close(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	filescanner := bufio.NewScanner(buf)
	retainTokens := make(map[string]struct{})
	args := utils.ConstructReplacements("ulds")

	for filescanner.Scan() {
		if filescanner.Text() != "" {
			retainTokens[filescanner.Text()] = struct{}{}
		}

		if err := filescanner.Err(); err != nil {
			CheckError(err)
		}
	}

	// Start a mutation loop
	var tokens sync.Map
	stdText := ""

	var wg sync.WaitGroup

	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		ngrams := utils.MakeToken(stdText)

		for _, token := range ngrams {
			if len(token) >= 4 {
				tokens.Store(token, struct{}{})
			}
		}

		stringWord := stdText

		wg.Add(1)
		go func(stringWord string) {
			defer wg.Done()

			// Create the retain mask
			mask := utils.CreateRetainMask(stringWord, retainTokens, args, doMultiByte, doNumberOfReplacements)

			// Use the retain mask in mutation
			tokens.Range(func(key, value interface{}) bool {
				newWord := utils.ReplaceWordByMask(stringWord, mask, key.(string), args, doNumberOfReplacements, doFuzzAmount)

				// Ensure results contain the retain tokens
				if newWord != "" {
					for value := range retainTokens {
						if strings.Contains(newWord, value) {
							fmt.Println(newWord)
						}
					}
				}
				return true
			})
		}(stringWord)
	}
	wg.Wait()
}

// CalculateEntropy calculates the entropy of the input strings and only prints
// those below the threshold
//
// Args:
//
//	stdIn (*bufio.Scanner): Buffer of standard input
//	threshold (string): Threshold to use for entropy
//	doDeHex (bool): If $HEX[...] text should be processed
//	verbose (bool): If verbose stdText should be printed about masks
//
// Returns:
//
//	None
func CalculateEntropy(stdIn *bufio.Scanner, threshold string, doDeHex bool, verbose bool) {
	if models.IsStringInt(threshold) == false {
		CheckError(errors.New("Invalid Chunk Size"))
	}
	thresholdInt, err := strconv.Atoi(threshold)
	CheckError(err)

	stdText := ""
	for stdIn.Scan() {

		if utils.TestHexInput(stdIn.Text()) == true && doDeHex == true {
			plaintext, err := utils.DehexPlaintext(stdIn.Text())
			if err != nil {
				stdText = ""
			}
			stdText = plaintext
		} else {
			stdText = stdIn.Text()
		}

		entropy := utils.TestEntropy(stdText)
		if entropy < thresholdInt && entropy != 0 {
			if verbose {
				fmt.Printf("%s:%d\n", stdText, entropy)
			} else {
				fmt.Printf("%s\n", stdText)
			}
		}
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

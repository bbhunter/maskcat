// Package that contains the primary logic for maskcat
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jakewnuk/maskcat/pkg/cmd"
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
		cmd.MatchMasks(stdIn, os.Args[2])
	case "sub":
		cmd.SubMasks(stdIn, os.Args[2])
	case "mutate":
		cmd.MutateMasks(stdIn, os.Args[2])
	case "tokens":
		cmd.GenerateTokens(stdIn, os.Args[2])
	case "partial":
		cmd.GeneratePartialMasks(stdIn, os.Args[2])
	case "remove":
		cmd.GeneratePartialRemoveMasks(stdIn, os.Args[2])
	default:
		cmd.GenerateMasks(stdIn)
	}
}

// printUsage prints usage information for the program
func printUsage() {
	fmt.Println("OPTIONS: match sub mutate tokens partial remove")
	fmt.Println("EXAMPLE: stdin | maskcat match <MASK-FILE>")
	fmt.Println("EXAMPLE: stdin | maskcat sub <TOKENS-FILE>")
	fmt.Println("EXAMPLE: stdin | maskcat mutate <CHUNK-SIZE>")
	fmt.Println("EXAMPLE: stdin | maskcat tokens <TOKEN-LEN> (99+ returns all)")
	fmt.Println("EXAMPLE: stdin | maskcat partial <MASK-CHARS>")
	fmt.Println("EXAMPLE: stdin | maskcat remove <MASK-CHARS>")
}

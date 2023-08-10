// Package main controls the primary logic for the application
//
// The package leans on /internal/cli to perform command line actions
// The application logic is stored within /pkg/*
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jakewnuk/maskcat/internal/cli"
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
		cli.MatchMasks(stdIn, os.Args[2])
	case "sub":
		cli.SubMasks(stdIn, os.Args[2])
	case "mutate":
		cli.MutateMasks(stdIn, os.Args[2])
	case "tokens":
		cli.GenerateTokens(stdIn, os.Args[2])
	case "partial":
		cli.GeneratePartialMasks(stdIn, os.Args[2])
	case "remove":
		cli.GeneratePartialRemoveMasks(stdIn, os.Args[2])
	default:
		cli.GenerateMasks(stdIn)
	}
}

// printUsage prints usage information for the application
func printUsage() {
	fmt.Println("OPTIONS: match sub mutate tokens partial remove")
	fmt.Println("EXAMPLE: stdin | maskcat match <MASK-FILE>")
	fmt.Println("EXAMPLE: stdin | maskcat sub <TOKENS-FILE>")
	fmt.Println("EXAMPLE: stdin | maskcat mutate <CHUNK-SIZE>")
	fmt.Println("EXAMPLE: stdin | maskcat tokens <TOKEN-LEN> (99+ returns all)")
	fmt.Println("EXAMPLE: stdin | maskcat partial <MASK-CHARS>")
	fmt.Println("EXAMPLE: stdin | maskcat remove <MASK-CHARS>")
}

// Package main controls the primary logic for the application
//
// The package leans on /internal/cli to perform command line actions
// The application logic is stored within /pkg/*
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/jakewnuk/maskcat/internal/cli"
)

var version = "0.0.1"

func main() {
	flagSet := flag.NewFlagSet("maskcat", flag.ExitOnError)
	doVerbose := flagSet.Bool("v", false, "Show verbose information about masks\nExample: maskcat [MODE] -v")
	doMultiByte := flagSet.Bool("m", false, "Process multibyte text (warning: slows processes)\nExample: maskcat [MODE] -m")
	doDeHex := flagSet.Bool("d", false, "Process $HEX[...] text (warning: slows processes)\nExample: maskcat [MODE] -d")
	doNumberOfReplacements := flagSet.Int("n", 1, "Max number of replacements to make per item (default: 1)\nExample: maskcat [MODE] -n 1")
	doFuzzAmount := flagSet.Int("f", 0, "Adds extra fuzz to the replacement functions\nExample: maskcat [MODE] -f 1")
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Options for maskcat (version %s):\n\n", version)
		flagSet.PrintDefaults()
		printUsage()
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(flagSet.Output(), "Options for maskcat (version %s):\n\n", version)
		flagSet.PrintDefaults()
		printUsage()
		os.Exit(0)
	}

	stdIn := bufio.NewScanner(os.Stdin)

	switch os.Args[1] {
	case "mask":
		flagSet.Parse(os.Args[2:])
		cli.GenerateMasks(stdIn, *doMultiByte, *doDeHex, *doVerbose)
	case "match":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.MatchMasks(stdIn, os.Args[2], *doMultiByte, *doDeHex)
	case "sub":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.SubMasks(stdIn, os.Args[2], *doMultiByte, *doDeHex, *doNumberOfReplacements, *doFuzzAmount)
	case "mutate":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.MutateMasks(stdIn, os.Args[2], *doMultiByte, *doDeHex, *doNumberOfReplacements, *doFuzzAmount)
	case "tokens":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.GenerateTokens(stdIn, os.Args[2], *doDeHex)
	case "partial":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.GeneratePartialMasks(stdIn, os.Args[2], *doDeHex)
	case "remove":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.GeneratePartialRemoveMasks(stdIn, os.Args[2], *doDeHex)
	case "retain":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.GenerateTokenRetainMasks(stdIn, os.Args[2], *doMultiByte, *doDeHex, *doNumberOfReplacements)
	case "splice":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.GenerateSpliceMutation(stdIn, os.Args[2], *doMultiByte, *doDeHex, *doNumberOfReplacements, *doFuzzAmount)
	case "filter":
		cli.CheckIfArgExists(2, os.Args)
		flagSet.Parse(os.Args[3:])
		cli.CalculateEntropy(stdIn, os.Args[2], *doMultiByte, *doVerbose)
	}
}

// printUsage prints usage information for the application
func printUsage() {
	fmt.Println(fmt.Sprintf("\nModes for maskcat (version %s):", version))
	fmt.Println("\n  mask\t\tCreates masks from text")
	fmt.Println("\t\tExample: stdin | maskcat mask [OPTIONS]")
	fmt.Println("\n  match\t\tMatches text to masks")
	fmt.Println("\t\tExample: stdin | maskcat match [MASK-FILE] [OPTIONS]")
	fmt.Println("\n  sub\t\tReplaces text with text from a file with masks")
	fmt.Println("\t\tExample: stdin | maskcat sub [TOKENS-FILE] [OPTIONS]")
	fmt.Println("\n  mutate\tMutates text by using chunking and token swapping")
	fmt.Println("\t\tExample: stdin | maskcat mutate [MIN-TOKEN-SIZE] [OPTIONS]")
	fmt.Println("\n  tokens\tSplits text into tokens and only print certain lengths (values over 99 allow all)")
	fmt.Println("\t\tExample: stdin | maskcat tokens [TOKEN-LEN] [OPTIONS]")
	fmt.Println("\n  partial\tPartially replaces characters with mask characters")
	fmt.Println("\t\tExample: stdin | maskcat partial [MASK-CHARS] [OPTIONS]")
	fmt.Println("\n  remove\tRemoves characters that match given mask characters")
	fmt.Println("\t\tExample: stdin | maskcat remove [MASK-CHARS] [OPTIONS]")
	fmt.Println("\n  retain\tCreates retain masks by keeping text from a file")
	fmt.Println("\t\tExample: stdin | maskcat retain [TOKENS-FILE] [OPTIONS]")
	fmt.Println("\n  splice\tMutates text by using retain masks and token swapping")
	fmt.Println("\t\tExample: stdin | maskcat splice [TOKENS-FILE] [OPTIONS]")
	fmt.Println("\n  filter\tOnly prints masks below a maximum entropy threshold")
	fmt.Println("\t\tExample: stdin | maskcat filter [ENTROPY-MAX] [OPTIONS]")
}

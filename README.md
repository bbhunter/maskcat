`Maskcat` is a multi-tool for working with text streams for password cracking.

Maskcat (`cat` mask) focuses on the usage of masks to extract and transform text and features several functions:

   - Making `hashcat` masks from `stdin`
   - Matching words from `stdin` to masks
   - Substituting tokens into `stdin` using masks
   - Mutating `stdin` with masks for new candidates
   - Generating tokens from `stdin` by extracting input
   - Creating partial masks from `stdin` by selecting character sets
   - Removing characters from `stdin` by selecting character sets
   - Creating retain masks from `stdin` by selecting tokens to retain

Maskcat also supports several options to assist in being a flexible and powerful tool:

- Multibyte text support
- Auto-dehexing text support
- Configurable number of replacements
- Additional fuzz configuration for replacements to create unique output

Maskcat fits into a small tool ecosystem for password cracking and is designed for lightweight and easy usage with its companion tools:

- [maskcat](https://github.com/JakeWnuk/maskcat)
- [rulecat](https://github.com/JakeWnuk/rulecat)
- [mode](https://github.com/JakeWnuk/mode)

### Getting Started

Usage information and other documentation can be found below:

- Usage documentation:
    - [Creating and Matching Masks](https://github.com/JakeWnuk/maskcat/blob/main/docs/CREATE_AND_MATCH.md)
    - [Token Swapping and Mutation](https://github.com/JakeWnuk/maskcat/blob/main/docs/SWAP_AND_MUTATE.md)
    - [Generating Tokens and Retain Masks](https://github.com/JakeWnuk/maskcat/blob/main/docs/TOKENS_AND_RETAIN.md)
    - [Partial Masks and Removing Character Sets](https://github.com/JakeWnuk/maskcat/blob/main/docs/PARTIAL_AND_REMOVE.md)

- For more application examples:
    - [Maskcat Examples](https://jakewnuk.com/posts/advanced-maskcat-cracking-guide/) (external link)

### Install from Source
```
git clone https://github.com/JakeWnuk/maskcat && cd maskcat && go build ./cmd/maskcat && mv ./maskcat ~/go/bin/
```

### Current Version 2.1.0:

```
Options for maskcat (version 2.1.0):

  -d    Process $HEX[...] text (warning: slows processes)
        Example: maskcat [MODE] -d
  -f int
        Adds extra fuzz to the replacement functions
        Example: maskcat [MODE] -f 1
  -m    Process multibyte text (warning: slows processes)
        Example: maskcat [MODE] -m
  -n int
        Max number of replacements to make per item (default: 1)
        Example: maskcat [MODE] -n 1 (default 1)
  -v    Show verbose information about masks
        Example: maskcat [MODE] -v

Modes for maskcat (version 2.1.0):

  mask          Creates masks from text
                Example: stdin | maskcat mask [OPTIONS]

  match         Matches text to masks
                Example: stdin | maskcat match [MASK-FILE] [OPTIONS]

  sub           Replaces text with text from a file with masks
                Example: stdin | maskcat sub [TOKENS-FILE] [OPTIONS]

  mutate        Mutates text by using chunking and token swapping
                Example: stdin | maskcat mutate [CHUNK-SIZE] [OPTIONS]

  tokens        Splits text into tokens and only print certain lengths (values over 99 allow all)
                Example: stdin | maskcat tokens [TOKEN-LEN] [OPTIONS]

  partial       Partially replaces characters with mask characters
                Example: stdin | maskcat partial [MASK-CHARS] [OPTIONS]

  remove        Removes characters that match given mask characters
                Example: stdin | maskcat remove [MASK-CHARS] [OPTIONS]

  retain        Creates retain masks by keeping text from a file
                Example: stdin | maskcat retain [TOKENS-FILE] [OPTIONS]
```

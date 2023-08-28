<h1 align="center">
Maskcat
 </h1>

- Maskcat (`cat` mask) performs seven (7) functions:
    - Makes Hashcat masks from stdin. Format is `MASK:LENGTH:COMPLEXITY:ENTROPY`.
    - Matches words from `stdin` to masks.
    - Substitutes tokens in wordlists using masks.
    - Mutates `stdin` using masks to create new candidates.
    - Generates tokens from `stdin` by removing non-alpha characters.
    - Partially replaces masks from `stdin` by selecting character sets.
    - Removes characters from `stdin` by selecting character sets.
- `maskcat` supports multibyte text
- For more application examples:
    - [Maskcat Examples](https://jakewnuk.com/posts/advanced-maskcat-cracking-guide/)
- See also [rulecat](https://github.com/JakeWnuk/rulecat)

## Getting Started

- [Install](#install)
- [Making Masks](#making-masks)
- [Matching Words to Masks](#matching-words-to-masks)
- [Substituting Tokens in Words with Masks](#substituting-tokens-in-words-with-masks)
- [Mutating Input](#mutating-input)
- [Generating Tokens](#generating-tokens)
- [Partial Masks](#partial-masks)
- [Removing Characters](#removing-characters)

### Install
#### Go
```
# known issue where 2.0.0 is not pulling
go install -v github.com/jakewnuk/maskcat/cmd/maskcat@latest
```
#### Source
```
git clone https://github.com/JakeWnuk/maskcat && cd maskcat && go build ./cmd/maskcat && mv ./maskcat ~/go/bin/

```
```
Options for maskcat (version 2.0.0):
  -m    Process multibyte text (warning: slows processes)
        Example: maskcat [MODE] -m
  -n int
        Max number of replacements to make per item (default: 1)
        Example: maskcat [MODE] -n 1 (default 1)
  -v    Show verbose information about masks
        Example: maskcat [MODE] -v

Modes for maskcat (version 2.0.0):

  mask          Creates masks from text
                Example: stdin | maskcat mask [OPTIONS]

  match         Matches text to masks
                Example: stdin | maskcat match [MASK-FILE] [OPTIONS]

  sub           Replaces text with text from a file with masks
                Example: stdin | maskcat sub [TOKENS-FILE] [OPTIONS]

  mutate        Mutates text by using chunking and token swapping
                Example: stdin | maskcat mutate [CHUNK-SIZE] [OPTIONS]

  tokens        Splits text into chunks by length (values over 99 allow all)
                Example: stdin | maskcat tokens [TOKEN-LEN] [OPTIONS]

  partial       Partially replaces characters with mask characters
                Example: stdin | maskcat partial [MASK-CHARS] [OPTIONS]

  remove        Removes characters that match given mask characters
                Example: stdin | maskcat remove [MASK-CHARS] [OPTIONS]
```

## Making Masks:
- Makes Hashcat masks from stdin. Format is `MASK:LENGTH:COMPLEXITY:ENTROPY`.
 ```
$ echo 'ThisISaT3ST123!' | maskcat mask -v
?u?l?l?l?u?u?l?u?d?u?u?d?d?d?s:15:4:333
 ```

 ```
$ head cracked.lst | maskcat mask -v
?l?l?l?l?l?l?l?l?l?l?l?l?d:13:2:322
?l?l?l?l?l?l?l?l?l?d:10:2:244
?l?l?l?l?l?l?l?l?l?l?l?l?d:13:2:322
?l?l?l?l?l?l?l?l?l?l?l?d:12:2:296
?l?l?l?l?l?l?l?l?l?l?l?l?l?l?d?d:16:2:384
?l?l?l?l?l?l?l?l?l?l?l?l?l?l?l?l?d?d:18:2:436
?u?l?l?l?l?u?l?l?l?l?l?d?d:13:3:306
?l?l?l?l?l?l?l?l?l?l?l?l?d?d:14:2:332
?l?l?l?l?l?l?l?l?l?l?l?l?d?d:14:2:332
 ```

 ```
$ head -n 100 cracked.lst | maskcat mask -v | cut -d ':' -f1 | sort | uniq -c | sort -rn
    8 ?u?l?l?l?l?l?l?l?d?d?d?d?s
    7 ?u?l?l?l?l?l?l?s?d?d?d?d
    6 ?u?l?l?l?l?l?l?d?d?d?d?s
    6 ?u?l?l?l?d?d?d?d?d?d?d?d
    5 ?u?l?l?l?l?l?l?l?d?d?d?d
    4 ?u?l?l?l?l?l?l?l?s?d?d?d
    4 ?u?l?l?l?l?l?l?l?l?s?d?d?d?d
    4 ?u?l?l?l?l?l?l?l?l?s?d?d?d
    4 ?u?l?l?l?l?l?l?l?l?l?d?s
    4 ?u?l?l?l?l?l?l?l?d?d
    4 ?u?l?l?l?l?l?l?d?d
```

## Matching Words to Masks:
- Matches words from `stdin` to masks.
 ```
$ cat masks.txt
?u?l?l?l?u?u?l?u?d?u?u?d?d?d?s

$ echo 'ThisISaT3ST123!' | maskcat match masks.txt
ThisISaT3ST123!
 ```

 ```
$ cat masks.txt
?u?l?l?l?u?u?l?u?d?u?u?d?d?d?s
?l?l?l?l

$ cat words.txt
ThisISaT3ST123!
test
bark
tree
Tree
Bark
NoMatch123

$ cat words.txt | maskcat match masks.txt
ThisISaT3ST123!
test
bark
tree
```

## Substituting Tokens in Words with Masks:
- Substitutes tokens in wordlists using masks.
```
# get a list of probable tokens

$ cat tokens.lst
Keywrd

# then take your favorite wordlist
$ cat words.lst
TheGreat123
TheGreats123
Thefats123
Greaty12345!!

# and sub matching masks with your token
$ cat test.lst | maskcat sub tokens.lst
TheKeywrd123
Keywrds123
Keywrd12345!!
 ```

## Mutating Input:
- Mutates `stdin` using masks to create new candidates.

### How Does Mutation Work?
- Mutation takes input from `stdin` then tokenizes it based on the length and valid values are added to an array.
- The array is then used in substitution mode to create new candidates.
- The results from the process are nondeterministic.
```
$ head -n 5 w.tmp | shuf | maskcat mutate 6 | sort -u
awesomawesome1
awesomear4
larrybear4
mathisear4
mathismathise1
ms.birdy8
ms.navit6

$ head -n 5 w.tmp | shuf | maskcat mutate 6 | sort -u
awesomawesome1
awesomear4
larrlitaf7
larrybear4
mathisawlitaf7
mathisear4
mathismathise1
ms.birdy8
ms.litaf7
ms.navit6
```

## Generating Tokens
- Generates tokens from `stdin` by parsing using several methods.
- Accepts an integer value to filter for token length based on the provided input.

### How does it work?
- Token generation replaces all digit and special characters within a string.
- Token generation find and splits camel casing.
- Token generation find and splits special and digit boundaries.
```
$ cat list.tmp
Password123
NotAPassword456

# Fetches all 8 length strings
$ cat list.tmp | maskcat tokens 8
Password
Password
Password

# If value is above 99 all tokens are allowed
$ cat list.tmp | maskcat tokens 99
Password
Password
Not
A
Password
NotAPassword
```

## Partial Masks
- Partially replaces masks from `stdin` by selecting character sets.
```
# Provide ulds as input and partial masks will be returned
$ cat list.tmp | maskcat partial d
Password?d?d?d
NotAPassword?d?d?d

# Multiple can also be used at once
$ cat list.tmp | maskcat partial du
?uassword?d?d?d
?uot?u?uassword?d?d?d
```

## Removing Characters
- Removes characters from `stdin` by selecting character sets.
```
# Provide ulds as input and remaining characters will be returned
$ echo 'Password123' | maskcat remove ul
123
```

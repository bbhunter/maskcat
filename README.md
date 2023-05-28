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
```
go install -v github.com/jakewnuk/maskcat@latest
```
```
OPTIONS: match sub mutate tokens partial remove
EXAMPLE: stdin | maskcat match <MASK-FILE>
EXAMPLE: stdin | maskcat sub <TOKENS-FILE>
EXAMPLE: stdin | maskcat mutate <CHUNK-SIZE>
EXAMPLE: stdin | maskcat tokens <TOKEN-LEN> (99+ returns all)
EXAMPLE: stdin | maskcat partial <MASK-CHARS>
EXAMPLE: stdin | maskcat remove <MASK-CHARS>
```

## Making Masks:
- Makes Hashcat masks from stdin. Format is `MASK:LENGTH:COMPLEXITY:ENTROPY`.
 ```
$ echo 'ThisISaT3ST123!' | maskcat
?u?l?l?l?u?u?l?u?d?u?u?d?d?d?s:15:4:333
 ```

 ```
$ head cracked.lst | maskcat 
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
$ head -n 100 cracked.lst | maskcat | cut -d ':' -f1 | sort | uniq -c | sort -rn
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
- Generates tokens from `stdin` by removing non-alpha characters.

### How does it work?
- Token generation replaces all digit and special characters within a string then filters for token length based on the provided input.
```
$ cat list.tmp
Password123
NotAPassword456

# Fetches all 8 length strings
$ echo 'Password123' | maskcat tokens 8
Password

# If value is above 99 all tokens are allowed
$ echo 'Password123' | maskcat tokens 99
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

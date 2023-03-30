<h1 align="center">
Maskcat
 </h1>

Maskcat performs 6 functions:
- Makes Hashcat masks from stdin. Format is `MASK:LENGTH:COMPLEXITY:ENTROPY`.
- Matches words from `stdin` to masks from a file argument.
- Substitutes tokens in wordlists using masks.
- Mutates `stdin` using masks to create new candidates.
- Generates tokens from `stdin` by removing non-alpha characters.
- Partially replaces masks from `stdin` by selecting character sets.

> NOTE: There is no support for `?b` or multi-byte characters at this time.

> For more examples: [link](https://jakewnuk.com/posts/advanced-maskcat-cracking-guide/)

## Getting Started

- [Install](#install)
- [Making Masks](#Making-Masks)
- [Matching Words to Masks](#Matching-Words-to-Masks)
- [Substituting Tokens in Words with Masks](#Substituting-Tokens-in-Words-with-Masks)
- [Mutating Input](#Mutating-Input)
- [Generating Tokens](#Generating-Tokens)
- [Partial Masks](#Partial-Masks)

### Install
```
go install -v github.com/jakewnuk/maskcat@latest
```

## Making Masks:

 ```
$ echo 'ThisISaT3ST123!' | maskcat
?u?l?l?l?u?u?l?u?d?u?u?d?d?d?s:15:4:333
 ```

 ```
$ head cracked.lst | maskcat 
 ?u?l?l?d?l?l?s?l?d?l?d?d?l:13:4
 ?d?l?d?l?u?d?u?u:8:3
 ?u?l?l?l?l?l?l?l?l?l?l?d:12:3
 ?u?l?d?d?l?l?l?l?l?l?l?l:12:3
 ?d?d?l?l?u?u?u?l:8:3
 ?u?l?l?l?l?l?l?s?l?l?l?l?l?d?d:15:4
 ?u?l?l?l?l?l?d?l?s?d?d:11:4
 ?u?l?l?l?d?d?d?s:8:4
 ?d?d?d?s?u?l?l?l?l?l?d?d?d:13:4
 ?u?l?l?l?l?u?l?l?d?d?s?d:12:4
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

### How Does Mutation Work?

- Mutation takes input from STDIN then tokenizes it based on the parameter
  provided by length. Tokens are then checked for length against the provided
  parameter and valid values are added to an array. This array is then used in
  the substitution mode to create new candidates. The results from the process
  are nondeterministic.

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

### How does it work?

- Token generation replaces all digit and special characters within a string then filters for token length based on the provided input.
- Spaces are preserved

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


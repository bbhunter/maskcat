### Quick Start
File examples used
```
$ cat test.txt
this
is a
big old test
luvTesting
it123always
works

$ cat retain.txt
this
test
works
123
```
Make tokens of length X
```
$ cat test.txt | maskcat tokens 3
isa
big
old
luv
```
Make retain masks from file input
```
$ cat test.txt | maskcat retain retain.txt
works
this
?l?l?s?l
?l?l?l?s?l?l?l?stest
?l?l123?l?l?l?l?l?l
?l?l?l?u?l?l?l?l?l?l
```

### Making Tokens
Maskcat can be used to create tokens of X length from `stdin` based on multiple
parsing methods. This will parse out tokens from input strings into smaller
alphabetical substrings that can be used with other modes.

This is used to identify trends and patterns in material for other use cases.

```
Example: stdin | maskcat tokens [TOKEN-LEN] [OPTIONS]
```

The `tokens` mode is affected by the following option flags:
- `-d` to process `$HEX[...]` text

When the `TOKEN-LEN` value is above 99 all tokens are allowed through. The
tokenizer can parse the following items:
 - Parses out camel case
 - Parses out digit boundaries
 - Parses out special characters boundaries
 - Parses out non-alpha characters

The following regex are used to do this:
 - `[A-Z][a-z]*|\d+|[^\dA-Z]+`
 - `[A-Z][a-z]*|\d+|\W+|\w+`
 - `[^a-zA-Z]+`

### Making Retain Masks
Maskcat can be used to create retain masks from `stdin` by creating masks
except for tokens given in a file. This will transform input into masks except
for any input given from a file.

This can be used to create masks that preserve common material.
```
Example: stdin | maskcat retain [TOKENS-FILE] [OPTIONS]
```

The `retain` mode is affected by the following option flags:
- `-m` to process multibyte text
- `-d` to process `$HEX[...]` text
- `-n` to control the max number of replacements per string

When the `-n` or max number of replacements value is provided the default (1)
number of max replacements can be changed.
```
$ cat retain.txt
this
test
works
123
old
always

$ cat test.txt | maskcat retain retain.txt -n 2
this
?l?l123always
?l?l?l?sold?stest
works
?l?l?l?u?l?l?l?l?l?l
?l?l?s?l
```

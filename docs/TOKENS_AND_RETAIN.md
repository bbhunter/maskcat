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
```
Make tokens of length X
```
$ cat test.txt | maskcat tokens 3
isa
big
old
luv
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


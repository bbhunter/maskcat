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
Filter and print masks with entropy value less than 100
```
$ cat test.txt | maskcat partial d | maskcat entropy 100 -v
it?d?d?dalways:10
```
```
$ cat test.txt | maskcat partial ds | maskcat entropy 100 -v
is?sa:5
big?sold?stest:10
it?d?d?dalways:10
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

### Filtering Masks by Entropy
Maskcat can be used to filter masks from `stdin` that are greater than a target
entropy value. This will only print items to `stdout` that are below the target
threshold. 

This is used to remove large masks from output that will not be feasible with
mask attacks. The `entropy` mode can use both full and partial masks.

```
Example: stdin | maskcat entropy [ENTROPY-MAX] [OPTIONS]
```

The `entropy` mode is affected by the following option flags:
- `-d` to process `$HEX[...]` text
- `-v` to print the entropy value of each item

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
Make partial masks by character set
```
$ cat test.txt | maskcat partial uds
this
is?sa
big?sold?stest
luv?uesting
it?d?d?dalways
works
```
Remove characters from tokens by character set
```
$ cat test.txt | maskcat remove uds
this
isa
bigoldtest
luvesting
italways
works
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

### Making Partial Masks
Maskcat can be used to create partial masks from `stdin` based on a provided
character set. This will replace any matching characters with their mask
equivalent to make a partial mask.

This is used to create partial masks and identify trends.

```
Example: stdin | maskcat partial [MASK-CHARS] [OPTIONS]
```

The `partial` mode is affected by the following option flags:
- `-d` to process `$HEX[...]` text

The following `MASK-CHARS` values are allowed:
- `u` to process upper case characters
- `l` to process lower case characters
- `d` to process digit characters
- `s` to process special characters
- `b` to process byte characters

### Removing Characters
Maskcat can be used to remove characters from `stdin` based on a provided
character set. This will remove any matching characters by their mask
equivalent to make a new string.

This is used to remove characters from strings for other use cases.
```
Example: stdin | maskcat remove [MASK-CHARS] [OPTIONS]
```

The `remove` mode is affected by the following option flags:
- `-d` to process `$HEX[...]` text

The following `MASK-CHARS` values are allowed:
- `u` to process upper case characters
- `l` to process lower case characters
- `d` to process digit characters
- `s` to process special characters
- `b` to process byte characters


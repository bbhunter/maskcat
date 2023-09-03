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


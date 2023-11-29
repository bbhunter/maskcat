### Quick Start
File examples used
```
$ cat test.txt
this
is a
Test!
love
Testing
Work$

$ cat retain.txt
Test
```
Create spliced text
```
$ cat test.txt | maskcat splice retain.txt
Test$
```
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

### Splicing Text
Maskcat can be used to "splice" text from `stdin` by parsing items from `stdin`
and inserting them into future items by using retain masks. This will transform
strings by shuffling tokens within them but also preserving content provided in
a file. 

This combines the `mutate` and `retain` mode into a single function to create 
token swapped text but ensuring certain tokens are not swapped and contained
within the results. 

The `splice` mode will use the tokenizer logic from the `tokens` mode to
generate substrings to use in the mutation logic.

```
Example: stdin | maskcat splice [TOKENS-FILE] [OPTIONS]
```

The `splice` mode is affected by the following option flags:
- `-m` to process multibyte text
- `-d` to process `$HEX[...]` text
- `-n` to control the max number of replacements per string
- `-f` to control the amount of extra fuzz to add to the replacements

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

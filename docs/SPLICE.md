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


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

$ cat sub.txt
Tree
Fruit
```
Substitute text from a file
```
$ cat test.txt | maskcat sub sub.txt
Tree$
Treeing
Tree!
Fruitng
```
Mutate text to create new candidates
```
$ cat test.txt | shuf | maskcat mutate 4
Work$
Test$
Work!
Test!
Tthisng
Working
Testing
Tloveng
```

### Swapping Text
Maskcat can be used to substitute text from `stdin` with items from a file
based on masks. This will transform the strings into masks and find masks where
the substring is allowed. If allowed a token swap is made and the new string is
created. 

This is often used to create new candidates based on material. Using `sub` mode
with `match` mode can help filter items.

```
Example: stdin | maskcat sub [TOKENS-FILE] [OPTIONS]
```

The `sub` mode is affected by the following option flags:
- `-m` to process multibyte text
- `-d` to process `$HEX[...]` text
- `-n` to control the max number of replacements per string
- `-f` to control the amount of extra fuzz to add to the replacements

When the `-n` flag is provided the default max number of replacements (1) can
be increased.
```
$ cat sub.txt
swap

$ echo 'string needs replacements' | maskcat sub sub.txt -n 2
swapng swaps replacements
```

When the `-f` flag is provided the default fuzz amount (0) can be increased.
The fuzzer works by extending the mask by its last value inside it. This allows
for tokens that normally would not match an opportunity to create a new
candidate.
```
$ cat sub.txt
stringz
bello

$ echo 'hello world' | maskcat sub sub.txt -f 2
hello stringz
bello world
```

### Mutating Text
Maskcat can be used to mutate text from `stdin` by parsing items from `stdin`
and inserting them into future items. This will transform strings by shuffling
tokens within them. If the token length is greater than or equal to the input length then it is
kept otherwise discarded and not used for swapping.

This is often used to rapidly create new candidates based on material.

```
Example: stdin | maskcat mutate [MIN-TOKEN-SIZE] [OPTIONS]
```

The `mutate` mode is affected by the following option flags:
- `-m` to process multibyte text
- `-d` to process `$HEX[...]` text
- `-n` to control the max number of replacements per string
- `-f` to control the amount of extra fuzz to add to the replacements

The `mutate` mode will use the tokenizer logic from the `tokens` mode to
generate substrings to use in the mutation logic. The `mutate` mode is non-deterministic by design
this means that it will not generate identical output for the same input.

When `mutate` is used the items are being read one-by-one from standard input and goes through the following steps at a high level:
- Parse text for ngrams/tokens and add them to the map
- Iterate over the map of tokens and perform swaps on possible tokens

For example, if "Test123" goes in and the ngram "Test" is pulled out and added to the map. Next, "Work345" is entered and the token "Work" is pulled out and added to the map.
When it comes time for the candidates to be made for "Work345" the map contains "Test" and "Work" making the possible options "Test345" and "Work345".
However, if the order was reversed and "Work345" was processed first and "Test123" was processed after then "Work123" and "Test123" would be the output.

Once the program is started, the map begins to fill with different items and depending on the order in which they are processed the output could be different.
This can also be multiplied by using the `shuf` command to mix up in the input and goroutines will also process items in a different order due to the multiple "threads" being used.

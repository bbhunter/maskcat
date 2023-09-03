### Quick Start
Create masks
```
$ echo 'This is a T3st!' | maskcat mask
?u?l?l?l?s?l?l?s?l?s?u?d?l?l?s
```

To show the `LENGTH:COMPLEXITY:ENTROPY` use the `-v` or `verbose` flag
```
$ echo 'This is a T3st!' | maskcat mask -v
?u?l?l?l?s?l?l?s?l?s?u?d?l?l?s:15:4:402
```

Match masks from a file
```
$ cat match.txt
?u?l?l?l?s?l?l?s?l?s?u?d?l?l?s

$ echo 'This is a T3st!' | maskcat match match.txt
This is a T3st!
```

### Creating Masks
Maskcat can be used to transform `stdin` into `Hashcat` masks using the `mask`
mode. This will replace all valid characters with their placeholder equivalent
and can then be used by other modes or tools.

```
Example: stdin | maskcat mask [OPTIONS]
```

The `mask` mode is affected by the following option flags:
- `-m` to process multibyte text
- `-d` to process `$HEX[...]` text
- `-v` to show verbose information about the mask

When the `-v` flag is provided the output format is:
- `MASK:LENGTH:COMPLEXITY:ENTROPY`

### Matching Masks
Maskcat can be used to match input from `stdin` to masks from a given file.
Matching items will be printed to `stdout` and this mode is often used to
filter or reduce wordlists to likely patterns.

```
Example: stdin | maskcat match [MASK-FILE] [OPTIONS]
```

The `match` mode is affected by the following option flags:
- `-m` to process multibyte text
- `-d` to process `$HEX[...]` text

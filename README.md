# `splitfiles`

```
usage: splitfiles [<flags>] <PATTERN> <TEMPLATE>

Splits STDIN into files when encountering a pattern.

The generated filenames will be printed to STDOUT. With the verbose flag, the number of
lines in each file will also be printed.

The pattern can be either a list of literal characters, or a regular expression that work
line by line. The line break character is part of the line you are matching by.

The output will not contain the split characters.

Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
  -f, --force    Overwrite files instead of skipping them.
  -E, --regexp   Parse PATTERN as a regular expression instead of a raw string.
  -v, --verbose  Print number of lines in each file.

Args:
  <PATTERN>   Pattern to split on.
  <TEMPLATE>  File template to generate from. You can control where in the filenames the
              sequential number will appear by inserting a series of "?" in it.
```

## Installation

If you have Go installed and a proper `$GOPATH` set up, you can install using the normal `go install` command:

```bash
go install github.com/Mange/splitfiles
```


### Autocompletion

You can then install shell autocompletion by adding one of these lines to your RC file:

```bash
# For bash
eval "$(splitfiles --completion-script-bash)"

# For ZSH
eval "$(splitfiles --completion-script-zsh)"
```

### Man page

You can install a man page with the following command:

```bash
splitfiles --help-man > /usr/local/share/man/man1/splitfiles.1
```

Depending on your machine, you might need elevated priviliges to write to this directoy.

## License

See LICENSE file, but tl;dr: MIT.

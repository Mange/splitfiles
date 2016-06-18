package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"

	"github.com/Mange/splitfiles/splitter"
)

const usageDescription = `Splits STDIN into files when encountering a pattern.

The generated filenames will be printed to STDOUT. With the verbose flag, the
number of lines in each file will also be printed.

The pattern can be either a list of literal characters, or a regular
expression that work line by line. The line break character is part of the line
you are matching by.

The output will not contain the split characters.`

var (
	app      = kingpin.New("splitfiles", usageDescription)
	pattern  = app.Arg("PATTERN", "Pattern to split on.").Required().String()
	template = app.Arg("TEMPLATE", "File template to generate from.\nYou can control where in the filenames the sequential number will appear by inserting a series of \"?\" in it.").Required().String()

	overwrite       = app.Flag("force", "Overwrite files instead of skipping them.").Short('f').Bool()
	patternIsRegexp = app.Flag("regexp", "Parse PATTERN as a regular expression instead of a raw string.").Short('E').Bool()
	verbose         = app.Flag("verbose", "Print number of lines in each file.").Short('v').Bool()
)

const (
	verboseOutput = "\t%d\n"
	normalOutput  = "\n"
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	err := SetupFilenameTemplate(*template)
	if err != nil {
		app.FatalUsage(err.Error())
	}

	splitter, err := splitter.New(*pattern, *patternIsRegexp)
	app.FatalIfError(err, "Could not parse PATTERN as Regexp: ")

	file, err := openNextFile()
	app.FatalIfError(err, "Could not create file: ")

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(file)
	linesWritten := 0

	err = scanChunks(scanner, splitter, func(chunk string, newChunk bool) error {
		if newChunk {
			writer.Flush()
			file.Close()
			printLineswritten(linesWritten)

			file, err = openNextFile()
			if err != nil {
				return err
			}
			writer = bufio.NewWriter(file)
			linesWritten = 0
		}

		_, err := fmt.Fprint(writer, chunk)

		// Special case: If we write any contents at all to a file, but that line
		// never ends with a newline, we should still count one line. It's not the
		// end of a line that defines a line, it's the start of it.
		// However, if we keep on writing empty chunks to the file and the file
		// ends up with zero bytes written, it has no lines.
		if linesWritten == 0 && len(chunk) > 0 {
			linesWritten += 1
		}

		linesWritten += strings.Count(chunk, "\n")

		return err
	})
	app.FatalIfError(err, "")

	writer.Flush()
	file.Close()
	printLineswritten(linesWritten)
}

func printLineswritten(linesWritten int) {
	if *verbose {
		fmt.Printf("\t%d\n", linesWritten)
	} else {
		fmt.Print("\n")
	}
}

func openNextFile() (*os.File, error) {
	filename := NextFilename()
	_, err := os.Stat(filename)
	exists := !os.IsNotExist(err)

	if exists && !*overwrite {
		app.Errorf("File %s already exists. Skipping it.", filename)
		return openNextFile()
	} else {
		fmt.Print(filename)
		return os.Create(filename)
	}
}

func scanChunks(
	scanner *bufio.Scanner,
	splitter splitter.Splitter,
	block func(string, bool) error,
) error {
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		parts := splitter.Split(line)

		if len(parts) == 1 {
			err := block(line, false)
			if err != nil {
				return err
			}
		} else {
			/*
				First part is not "new", but all the others are:

				foo bar baz, split on " ":
					"foo" is the end of the last chunk
					"bar" is a new chunk, which also ends here
					"baz" is a new chunk, which continues
			*/
			for i := 0; i < len(parts); i++ {
				block(parts[i], i > 0)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
		app.Fatalf("Error while reading STDIN: %s", err.Error())
	}
	return nil
}

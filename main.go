package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"strings"
)

var (
	app      = kingpin.New("splitfiles", "Splits STDIN into files when encountering a pattern.")
	pattern  = app.Arg("pattern", "Pattern to split on.").Required().String()
	template = app.Arg("template", "File template to generate from.\nYou can control where in the filenames the sequential number will appear by inserting a series of \"?\" in it.").Required().String()

	overwrite = app.Flag("force", "Overwrite files instead of skipping them").Short('f').Bool()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	err := SetupFilenameTemplate(*template)
	if err != nil {
		app.FatalUsage(err.Error())
	}

	file, err := openNextFile()
	app.FatalIfError(err, "Could not create file: ")

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(file)

	err = scanLines(scanner, func(line string, newChunk bool) error {
		if newChunk {
			writer.Flush()
			file.Close()

			file, err = openNextFile()
			if err != nil {
				return err
			}
			writer = bufio.NewWriter(file)
		}

		_, err := fmt.Fprintln(writer, line)
		return err
	})
	app.FatalIfError(err, "")

	writer.Flush()
	file.Close()
}

func openNextFile() (*os.File, error) {
	filename := NextFilename()
	_, err := os.Stat(filename)
	exists := !os.IsNotExist(err)

	if exists && !*overwrite {
		app.Errorf("File %s already exists. Skipping it.", filename)
		return openNextFile()
	} else {
		return os.Create(filename)
	}
}

func scanLines(scanner *bufio.Scanner, block func(string, bool) error) error {
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, *pattern)

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

func handleLine(line string, writer io.Writer) (splitHappened bool, remainder string) {
	parts := strings.SplitN(line, *pattern, 2)
	if len(parts) == 1 {
		fmt.Fprintln(writer, line)
		return false, ""
	} else {
		fmt.Fprintln(writer, parts[0])
		return true, parts[1]
	}
}

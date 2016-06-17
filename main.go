package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
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

	file, err := OpenNextFile()
	app.FatalIfError(err, "Could not create file: ")

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(file)

	for scanner.Scan() {
		fmt.Fprintln(writer, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		app.Fatalf("Error while reading STDIN: %s", err.Error())
	}

	writer.Flush()
	file.Close()
}

func OpenNextFile() (*os.File, error) {
	filename := NextFilename()
	_, err := os.Stat(filename)
	exists := !os.IsNotExist(err)

	if exists && !*overwrite {
		app.Errorf("File %s already exists. Skipping it.", filename)
		return OpenNextFile()
	} else {
		return os.Create(filename)
	}
}

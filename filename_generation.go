package main

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	placeholderPattern = regexp.MustCompile(`\?+`)

	currentIndex = 0

	formatTemplate string
)

func SetupFilenameTemplate(template string) error {
	matches := placeholderPattern.FindAllString(template, 2)

	if matches == nil {
		return SetupFilenameTemplate(template + ".?")
	} else if len(matches) > 1 {
		return errors.New("Template contained more than 1 series of question marks.\nYou can only use question marks at one place in the template: \"hello_??.txt\" (not \"hello_?_?.txt\").")
	} else {
		formatTemplate = placeholderPattern.ReplaceAllStringFunc(template, func(match string) string {
			return fmt.Sprintf("%%0%dd", len(match))
		})
		return nil
	}
}

func NextFilename() string {
	currentIndex += 1
	return generateFilename(currentIndex)
}

func generateFilename(index int) string {
	return fmt.Sprintf(formatTemplate, index)
}

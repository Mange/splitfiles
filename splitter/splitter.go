package splitter

import "regexp"
import "strings"

type Splitter interface {
	Split(string) []string
}

type regexpSplitter struct {
	pattern *regexp.Regexp
}

func newRegexpSplitter(pattern string) (*regexpSplitter, error) {
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &regexpSplitter{pattern: matcher}, nil
}

func (self *regexpSplitter) Split(text string) []string {
	return self.pattern.Split(text, -1)
}

type characterSplitter struct {
	characters string
}

func newCharacterSplitter(pattern string) *characterSplitter {
	return &characterSplitter{characters: pattern}
}

func (self *characterSplitter) Split(text string) []string {
	return strings.Split(text, self.characters)
}

func New(pattern string, patternIsRegexp bool) (Splitter, error) {
	if patternIsRegexp {
		return newRegexpSplitter(pattern)
	} else {
		return newCharacterSplitter(pattern), nil
	}
}

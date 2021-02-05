package main

import (
	"regexp"
)

// LineGrep is...
type LineGrep struct {
	re *regexp.Regexp
}

func (instance *LineGrep) init(pattern string, ignoreCase bool) (*LineGrep, error) {
	var err error

	if ignoreCase {
		pattern = "(?i)" + pattern
	}

	instance.re, err = regexp.Compile(pattern)
	return instance, err
}

func (instance *LineGrep) search(line string) []int {
	return instance.re.FindStringIndex(line)
}

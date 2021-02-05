package main

import (
	"regexp"
)

// LineGrep is...
type LineGrep struct {
	re *regexp.Regexp
}

func (instance *LineGrep) init(pattern string) (err error) {
	instance.re, err = regexp.Compile(pattern)
	return
}

func search(line string) (uint32, uint32, error) {
	return 0, 0, nil
}

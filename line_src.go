package main

import (
	"bufio"
	"io"
	"strings"
)

// LineSource is...
type LineSource interface {
	produce() (string, uint32, string, error)
}

// LinesFromStdin is...
type LinesFromStdin struct {
	reader *bufio.Scanner
	index  uint32
}

func (source LinesFromStdin) produce() (string, uint32, string, error) {
	line := ""
	prefix := ""
	var err error = nil

	if source.reader.Scan() {
		line = strings.Trim(source.reader.Text(), "\r\n")
	}

	if len(line) == 0 {
		err = io.EOF
	} else {
		err = source.reader.Err()
	}

	return prefix, source.index, line, err
}

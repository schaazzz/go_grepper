package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func (source *LinesFromStdin) init(reader *bufio.Scanner) *LinesFromStdin {
	source.reader = reader
	source.index = 0
	return source
}

func (source *LinesFromStdin) produce() (string, uint32, string, error) {
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

// LinesFromFiles is...
type LinesFromFiles struct {
	files           []string
	reader          *bufio.Scanner
	emitFilename    bool
	currentFilename string
	index           uint32
}

func (source *LinesFromFiles) init(files []string) *LinesFromFiles {
	source.files = files
	source.emitFilename = len(files) > 1
	source.currentFilename = ""
	source.index = 0
	return source
}

func (source *LinesFromFiles) produce() (string, uint32, string, error) {
	var line string
	var err error

	for {
		if nil != source.reader {
			if source.reader.Scan() {
				source.index++
				err = nil
				line = source.reader.Text()
				break
			} else {
				source.reader = nil
			}
		} else {
			var path string
			if len(source.files) == 0 {
				source.currentFilename = ""
				source.index = 0
				line = ""
				err = io.EOF
				break
			}

			path, source.files = source.files[0], source.files[1:]
			fileInfo, err := os.Stat(path)

			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: \"%s\" - no such file!", path)
				continue
			}

			if fileInfo.IsDir() {
				fmt.Fprintf(os.Stderr, "Info: \"%s\" is a directory!", path)
			}

			if source.currentFilename = ""; source.emitFilename {
				_, source.currentFilename = filepath.Split(path)
			}

			file, err := os.Open(path)
			if err == nil {
				source.index = 0
				source.reader = bufio.NewScanner(file)
			} else {
				fmt.Fprintf(os.Stderr, "Error: Can't open \"%s\"!", path)
			}
		}
	}

	return source.currentFilename, source.index, line, err
}

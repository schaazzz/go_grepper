package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"rsc.io/getopt"
)

// Flags is...
type Flags struct {
	ignoreCase    bool
	printLineNums bool
	useColor      bool
}

// Config is...
type Config struct {
	flags   Flags
	pattern string
	files   []string
}

// ResultIndex is...
type ResultIndex struct {
	start int
	end   int
}

func parseArgs() Config {
	var config Config

	flag.BoolVar(&config.flags.ignoreCase, "i", false, "case distinctions in patterns and data")
	flag.BoolVar(&config.flags.printLineNums, "n", false, "print line number with output lines")
	flag.BoolVar(&config.flags.useColor, "color", false, "use markers to highlight the matching strings")
	help := flag.Bool("help", false, "Print this help menu")

	getopt.Parse()
	if *help {
		flag.Usage()
	}

	config.pattern = flag.Args()[0]
	config.files = append(config.files, flag.Args()[1:len(flag.Args())]...)

	if len(config.files) > 0 {
		var files []string
		for _, item := range config.files {
			matches, _ := filepath.Glob(item)
			files = append(files, matches...)
		}
		config.files = files
	}

	return config
}

func process(source LineSource, flags *Flags, pattern string) {
	grep, err := new(LineGrep).init(pattern, flags.ignoreCase)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid regular expression!")
		return
	}

	exit := make(chan bool)
	mainTx := make(chan string)
	mainRx := make(chan ResultIndex)
	var wg sync.WaitGroup

	wg.Add(1)
	go func(rx chan string, tx chan ResultIndex, exit chan bool) {
		defer wg.Done()

	loop:
		for {
			select {
			case line := <-rx:
				result := grep.search(line)
				if len(result) == 2 {
					tx <- ResultIndex{start: result[0], end: result[1]}
				} else {
					tx <- ResultIndex{start: 0, end: 0}
				}
			case <-exit:
				break loop
			default:
				continue
			}
		}
	}(mainTx, mainRx, exit)

	for {
		prefix, index, line, err := source.produce()

		if err == io.EOF {
			exit <- true
			break
		}

		if err == nil {
			mainTx <- line
			result := <-mainRx

			if result.start != result.end {
				printMatchedLine(flags, prefix, index, line, result.start, result.end)
			}

		} else {
			fmt.Println("Error:", err)
		}
	}

	wg.Wait()
}

func printMatchedLine(flags *Flags, prefix string, index uint32, line string, start int, end int) (err error) {
	prefixWriter := color.New(color.FgMagenta)
	lineNumWriter := color.New(color.FgGreen)
	matchWriter := color.New(color.FgRed)
	semicolonWriter := color.New(color.FgCyan)

	color.NoColor = !flags.useColor
	if len(prefix) > 0 {
		prefixWriter.Print(prefix)
		semicolonWriter.Print(":")
	}

	if flags.printLineNums {
		lineNumWriter.Print(index)
		semicolonWriter.Print(":")
	}

	fmt.Print(line[:start])
	matchWriter.Print(line[start:end])
	fmt.Println(line[end:])

	return nil
}

func main() {
	config := parseArgs()

	reader := bufio.NewScanner(os.Stdin)
	if len(config.files) == 0 {
		process(
			new(LinesFromStdin).init(reader),
			&config.flags,
			config.pattern)
	} else {
		process(
			new(LinesFromFiles).init(config.files),
			&config.flags,
			config.pattern)
	}

	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
}

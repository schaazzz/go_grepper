package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

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

	return config
}

func process(source LineSource) {
	for {
		prefix, index, line, err := source.produce()

		if err == io.EOF {
			break
		}

		if err == nil {
			fmt.Printf("%s:%d:%s\n", prefix, index, line)
		} else {
			fmt.Println("Error:", err)
		}

	}
}

func main() {
	config := parseArgs()

	reader := bufio.NewScanner(os.Stdin)
	if len(config.files) == 0 {
		process(LinesFromStdin{reader: reader, index: 0})
	}

	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	// for reader.Scan() {
	// 	fmt.Println("!", reader.Text())
	// }

	// if err := reader.Err(); err != nil {
	// 	log.Println(err)
	// }
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type CountOptions struct {
	ByteCount bool
	LineCount bool
	WordCount bool
	Order     []string
}

func main() {
	options, filenames := parseArgs(os.Args[1:])

	if len(filenames) == 0 {
		fmt.Println("Usage: mwc [options] <filename>")
		fmt.Println("Options:")
		fmt.Println("  -l    Count Lines")
		fmt.Println("  -w    Count Words")
		fmt.Println("  -c    Count Bytes")
		os.Exit(1)
	}

	for _, filename := range filenames {
		counts, err := countFile(filename, options)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
			continue
		}
		printCounts(counts, filename, options.Order)
	}
}

func countFile(filename string, options CountOptions) (map[string]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	counts := make(map[string]int64)
	reader := bufio.NewReader(file)

	var byteCount, lineCount, wordCount int64
	inWord := false

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		byteCount++
		if b == '\n' {
			lineCount++
		}
		if unicode.IsSpace(rune(b)) {
			inWord = false
		} else {
			if !inWord {
				wordCount++
				inWord = true
			}
		}
	}

	if options.ByteCount {
		counts["bytes"] = byteCount
	}

	if options.LineCount {
		counts["lines"] = lineCount
	}

	if options.WordCount {
		counts["words"] = wordCount
	}

	return counts, nil
}

func printCounts(counts map[string]int64, filename string, order []string) {
	for _, countType := range order {
		if count, ok := counts[countType]; ok {
			fmt.Printf("   %d ", count)
		}
	}
	fmt.Println(filename)
}

func parseArgs(args []string) (CountOptions, []string) {
	options := CountOptions{}
	var filenames []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, char := range arg[1:] {
				switch char {
				case 'l':
					options.LineCount = true
					options.Order = append(options.Order, "lines")
				case 'w':
					options.WordCount = true
					options.Order = append(options.Order, "words")
				case 'c':
					options.ByteCount = true
					options.Order = append(options.Order, "bytes")
				}
			}
		} else {
			filenames = append(filenames, arg)
		}
	}

	return options, filenames
}

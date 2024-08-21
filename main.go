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
	ByteCount      bool
	LineCount      bool
	WordCount      bool
	CharacterCount bool
	Order          []string
}

func main() {
	options, filenames := parseArgs(os.Args[1:])

	if len(os.Args) == 1 || (len(filenames) == 0 && !hasAnyOption(options)) {
		printUsage()
		os.Exit(1)
	}

	if len(filenames) == 0 {
		// Read from stdin if no filename is provided
		counts, err := processInput(os.Stdin, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing stdin: %v\n", err)
			os.Exit(1)
		}
		printCounts(counts, "", options.Order)
	} else {
		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", filename, err)
				continue
			}
			counts, err := processInput(file, options)
			file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
				continue
			}
			printCounts(counts, filename, options.Order)
		}
	}
}

func processInput(input io.Reader, options CountOptions) (map[string]int64, error) {
	counts := make(map[string]int64)
	reader := bufio.NewReader(input)

	var byteCount, lineCount, wordCount, characterCount int64
	inWord := false

	for {
		r, size, err := reader.ReadRune() // Reads a single Unicode character (rune) from the input
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		// For ASCII text (where each character is one byte), byte count and character count will be the same.
		// For text with multibyte Unicode characters (like emoji or non-Latin scripts),
		//  byte count will be larger than character count.
		byteCount += int64(size)
		characterCount++

		if r == '\n' {
			lineCount++
		}
		if unicode.IsSpace(r) {
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

	if options.CharacterCount {
		counts["characters"] = characterCount
	}

	return counts, nil
}

func printCounts(counts map[string]int64, filename string, order []string) {
	for _, countType := range order {
		if count, ok := counts[countType]; ok {
			fmt.Printf("%8d", count)
		}
	}
	if filename != "" {
		fmt.Printf(" %s", filename)
	}
	fmt.Println()
}

func parseArgs(args []string) (CountOptions, []string) {
	options := CountOptions{}
	var filenames []string
	hasOptions := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			hasOptions = true
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
				case 'm':
					options.CharacterCount = true
					options.Order = append(options.Order, "characters")
				}
			}
		} else {
			filenames = append(filenames, arg)
		}
	}

	// If no options were provided, use the default options
	if !hasOptions {
		options.LineCount = true
		options.WordCount = true
		options.ByteCount = true
		options.Order = []string{"lines", "words", "bytes"}
	}

	return options, filenames
}

func printUsage() {
	fmt.Println("usage: mwc [-lwcm] [filename ...]")
	fmt.Println("Options:")
	fmt.Println("  -l    Count Lines")
	fmt.Println("  -w    Count Words")
	fmt.Println("  -c    Count Bytes")
	fmt.Println("  -m    Count Characters")
	fmt.Println("If no filename is provided, mwc reads from standard input.")
}

func hasAnyOption(options CountOptions) bool {
	return options.LineCount || options.WordCount || options.ByteCount || options.CharacterCount
}

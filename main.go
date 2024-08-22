//go:build mwc
// +build mwc

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// CountOptions holds the flags for different counting options
type CountOptions struct {
	ByteCount      bool
	LineCount      bool
	WordCount      bool
	CharacterCount bool
	Order          []string // Keeps track of the order in which options were specified
}

// FileCount holds the counts for a specific file
type FileCount struct {
	Filename string
	Counts   map[string]int64
}

func main() {
	options, filenames, err := parseArgs(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [-clmw] [file ...]\n", os.Args[0])
		os.Exit(1)
	}

	// If no arguments provided or no valid options, print usage and exit
	if len(os.Args) == 1 || (len(filenames) == 0 && !hasAnyOption(options)) {
		printUsage()
		os.Exit(1)
	}

	if len(filenames) == 0 {
		// Read from stdin if no filename is provided
		counts, err := processInput(os.Stdin, options)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error processing stdin: %v\n", err)
			os.Exit(1)
		}
		printCounts(counts, "", options.Order)
	} else {
		var fileCounts []FileCount
		totalCounts := make(map[string]int64)
		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", filename, err)
				continue
			}
			counts, err := processInput(file, options)
			_ = file.Close()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
				continue
			}
			fileCounts = append(fileCounts, FileCount{Filename: filename, Counts: counts})
			for k, v := range counts {
				totalCounts[k] += v
			}
		}

		// Print counts for each file
		for _, fc := range fileCounts {
			printCounts(fc.Counts, fc.Filename, options.Order)
		}

		// Print total if there's more than one file
		if len(fileCounts) > 1 {
			printCounts(totalCounts, "total", options.Order)
		}
	}
}

// processInput reads from the input and counts bytes, lines, words, and characters based on the options
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

	// Add counts to the map based on the options
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

// printCounts outputs the counts in the specified order
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

// parseArgs processes command-line arguments and returns CountOptions and filenames
func parseArgs(args []string) (CountOptions, []string, error) {
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
				default:
					//_, _ = fmt.Fprintf(os.Stderr, "%s: illegal option -- %c\n", os.Args[0], char)
					//_, _ = fmt.Fprintf(os.Stderr, "usage: %s [-clmw] [file ...]\n", os.Args[0])
					//os.Exit(1)
					return CountOptions{}, nil, fmt.Errorf("illegal option -- %c", char)
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

	return options, filenames, nil
}

// printUsage displays the usage information for the command
func printUsage() {
	fmt.Println("usage: mwc [-lwcm] [filename ...]")
	fmt.Println("Options:")
	fmt.Println("  -l    Count Lines")
	fmt.Println("  -w    Count Words")
	fmt.Println("  -c    Count Bytes")
	fmt.Println("  -m    Count Characters")
	fmt.Println("If no filename is provided, mwc reads from standard input.")
}

// hasAnyOption checks if any counting option is enabled
func hasAnyOption(options CountOptions) bool {
	return options.LineCount || options.WordCount || options.ByteCount || options.CharacterCount
}

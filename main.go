package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type CountOptions struct {
	ByteCount bool
}

func main() {
	options := parseFlag()

	if flag.NArg() < 1 {
		fmt.Println("Usage: mwc [options] <filename>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, filename := range flag.Args() {
		counts, err := countFile(filename, options)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
			continue
		}
		printCounts(counts, filename)
	}
}

func countFile(filename string, options CountOptions) (map[string]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	counts := make(map[string]int64)

	if options.ByteCount {
		reader := bufio.NewReader(file)
		var byteCount int64
		for {
			_, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error reading file: %w", err)
			}
			byteCount++
		}
		counts["bytes"] = byteCount
	}

	return counts, nil
}

func printCounts(counts map[string]int64, filename string) {
	for _, countType := range []string{"bytes"} {
		if count, ok := counts[countType]; ok {
			//fmt.Printf("%d %s ", count, countType)
			fmt.Printf("	%d ", count)
		}
	}
	fmt.Println(filename)
}

func parseFlag() CountOptions {
	byteCount := flag.Bool("c", false, "Count Bytes")

	flag.Parse()
	return CountOptions{
		ByteCount: *byteCount,
	}
}

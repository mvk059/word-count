package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestProcessInput tests all counting options with multiple inputs
func TestProcessInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  CountOptions
		expected map[string]int64
	}{
		{
			name:     "Byte Count - ASCII",
			input:    "Hello, World!\n",
			options:  CountOptions{ByteCount: true, Order: []string{"bytes"}},
			expected: map[string]int64{"bytes": 14},
		},
		{
			name:     "Byte Count - Unicode",
			input:    "Hello, 世界!\n",
			options:  CountOptions{ByteCount: true, Order: []string{"bytes"}},
			expected: map[string]int64{"bytes": 15},
		},
		{
			name:     "Line Count - Single Line",
			input:    "Hello, World!",
			options:  CountOptions{LineCount: true, Order: []string{"lines"}},
			expected: map[string]int64{"lines": 0},
		},
		{
			name:     "Line Count - Multiple Lines",
			input:    "Hello, World!\nGoodbye, World!\n",
			options:  CountOptions{LineCount: true, Order: []string{"lines"}},
			expected: map[string]int64{"lines": 2},
		},
		{
			name:     "Word Count - Single Word",
			input:    "Hello",
			options:  CountOptions{WordCount: true, Order: []string{"words"}},
			expected: map[string]int64{"words": 1},
		},
		{
			name:     "Word Count - Multiple Words",
			input:    "Hello, World!\nGoodbye, World!\n",
			options:  CountOptions{WordCount: true, Order: []string{"words"}},
			expected: map[string]int64{"words": 4},
		},
		{
			name:     "Character Count - ASCII",
			input:    "Hello, World!\n",
			options:  CountOptions{CharacterCount: true, Order: []string{"characters"}},
			expected: map[string]int64{"characters": 14},
		},
		{
			name:     "Character Count - Unicode",
			input:    "Hello, 世界!\n",
			options:  CountOptions{CharacterCount: true, Order: []string{"characters"}},
			expected: map[string]int64{"characters": 11},
		},
		{
			name:    "Default Option",
			input:   "Hello, World!\nGoodbye, World!\n",
			options: CountOptions{LineCount: true, WordCount: true, ByteCount: true, Order: []string{"lines", "words", "bytes"}},
			expected: map[string]int64{
				"lines": 2,
				"words": 4,
				"bytes": 30,
			},
		},
		{
			name:    "All Options",
			input:   "Hello, 世界!\nGoodbye, World!\n",
			options: CountOptions{LineCount: true, WordCount: true, ByteCount: true, CharacterCount: true, Order: []string{"lines", "words", "bytes", "characters"}},
			expected: map[string]int64{
				"lines":      2,
				"words":      4,
				"bytes":      31,
				"characters": 27,
			},
		},
		{
			name:     "Empty Input",
			input:    "",
			options:  CountOptions{LineCount: true, WordCount: true, ByteCount: true, CharacterCount: true, Order: []string{"lines", "words", "bytes", "characters"}},
			expected: map[string]int64{"lines": 0, "words": 0, "bytes": 0, "characters": 0},
		},
		{
			name:     "Only Whitespace",
			input:    "   \n\t\n  ",
			options:  CountOptions{LineCount: true, WordCount: true, ByteCount: true, CharacterCount: true, Order: []string{"lines", "words", "bytes", "characters"}},
			expected: map[string]int64{"lines": 2, "words": 0, "bytes": 8, "characters": 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			counts, err := processInput(input, tt.options)
			if err != nil {
				t.Fatalf("Error processing input: %v", err)
			}
			for k, v := range tt.expected {
				if counts[k] != v {
					t.Errorf("Expected %s: %d, got: %d", k, v, counts[k])
				}
			}
		})
	}
}

// TestStandardInput simulates reading from standard input
func TestStandardInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		args     []string
		expected map[string]int64
		order    []string
	}{
		{
			name:     "Standard Input - Single Line",
			input:    "Hello, World!",
			args:     []string{"-l"},
			expected: map[string]int64{"lines": 0},
			order:    []string{"lines"},
		},
		{
			name:     "Standard Input - Single Line with Newline",
			input:    "Hello, World!\n",
			args:     []string{"-l"},
			expected: map[string]int64{"lines": 1},
			order:    []string{"lines"},
		},
		{
			name:     "Standard Input - Multiple Lines",
			input:    "Hello, World!\nGoodbye, World!",
			args:     []string{"-l"},
			expected: map[string]int64{"lines": 1},
			order:    []string{"lines"},
		},
		{
			name:     "Standard Input - Multiple Lines with Newline",
			input:    "Hello, World!\nGoodbye, World!\n",
			args:     []string{"-l"},
			expected: map[string]int64{"lines": 2},
			order:    []string{"lines"},
		},
		{
			name:     "Standard Input - Default Options",
			input:    "Hello, World!\nGoodbye, World!\n",
			args:     []string{},
			expected: map[string]int64{"lines": 2, "words": 4, "bytes": 30},
			order:    []string{"lines", "words", "bytes"},
		},
		{
			name:     "Standard Input - All Options",
			input:    "Hello, World!\nGoodbye, World!\n",
			args:     []string{"-lwcm"},
			expected: map[string]int64{"lines": 2, "words": 4, "bytes": 30, "characters": 30},
			order:    []string{"lines", "words", "bytes", "characters"},
		},
		{
			name:     "Standard Input - Custom Order",
			input:    "Hello, World!\nGoodbye, World!\n",
			args:     []string{"-wlc"},
			expected: map[string]int64{"words": 4, "lines": 2, "bytes": 30},
			order:    []string{"words", "lines", "bytes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to simulate stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Error creating pipe: %v", err)
			}

			// Save the original stdin and args
			oldStdin := os.Stdin
			oldArgs := os.Args

			// Replace stdin with our pipe and set args
			os.Stdin = r
			os.Args = append([]string{"mwc"}, tt.args...)

			// Write the test input to the pipe
			go func() {
				defer w.Close()
				_, _ = w.Write([]byte(tt.input))
			}()

			// Capture stdout
			oldStdout := os.Stdout
			r2, w2, _ := os.Pipe()
			os.Stdout = w2

			// Run main
			main()

			// Restore stdout and close the write end of the pipe
			w2.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r2)
			output := strings.TrimSpace(buf.String())

			// Restore the original stdin and args
			os.Stdin = oldStdin
			os.Args = oldArgs

			// Parse the output
			fields := strings.Fields(output)
			if len(fields) != len(tt.expected) {
				t.Fatalf("Expected %d fields, got %d", len(tt.expected), len(fields))
			}

			// Check the counts
			for i, countType := range tt.order {
				expectedCount := tt.expected[countType]
				actualCount, err := strconv.ParseInt(fields[i], 10, 64)
				if err != nil {
					t.Fatalf("Error parsing count for %s: %v", countType, err)
				}
				if actualCount != expectedCount {
					t.Errorf("Expected %s count: %d, got: %d", countType, expectedCount, actualCount)
				}
			}
		})
	}
}

// TestMultipleFiles tests processing multiple files
func TestMultipleFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected []struct {
			lines int
			words int
			bytes int
			file  string
		}
	}{
		{
			name: "Two Files",
			files: map[string]string{
				"file1.txt": "Hello, World!\n",
				"file2.txt": "Goodbye, World!\n",
			},
			expected: []struct {
				lines int
				words int
				bytes int
				file  string
			}{
				{1, 2, 14, "file1.txt"},
				{1, 2, 16, "file2.txt"},
				{2, 4, 30, "total"},
			},
		},
		{
			name: "Three Files",
			files: map[string]string{
				"file1.txt": "Hello, World!\n",
				"file2.txt": "Goodbye, World!\n",
				"file3.txt": "Test file.\n",
			},
			expected: []struct {
				lines int
				words int
				bytes int
				file  string
			}{
				{1, 2, 14, "file1.txt"},
				{1, 2, 16, "file2.txt"},
				{1, 2, 11, "file3.txt"},
				{3, 6, 41, "total"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for test files
			tmpDir, err := os.MkdirTemp("", "wc_test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			var filenames []string
			for filename, content := range tt.files {
				path := filepath.Join(tmpDir, filename)
				err := os.WriteFile(path, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test file %s: %v", filename, err)
				}
				filenames = append(filenames, path)
			}

			// Sort filenames to ensure consistent order
			sort.Strings(filenames)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run main with test files
			os.Args = append([]string{"mwc"}, filenames...)
			main()

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := strings.TrimSpace(buf.String())
			lines := strings.Split(output, "\n")

			// Check output
			if len(lines) != len(tt.expected) {
				t.Fatalf("Expected %d lines of output, got %d", len(tt.expected), len(lines))
			}
			for i, expected := range tt.expected {
				parts := strings.Fields(lines[i])
				if len(parts) != 4 {
					t.Errorf("Line %d: expected 4 fields, got %d", i+1, len(parts))
					continue
				}

				actualLines, _ := strconv.Atoi(parts[0])
				actualWords, _ := strconv.Atoi(parts[1])
				actualBytes, _ := strconv.Atoi(parts[2])
				actualFile := filepath.Base(parts[3])

				if actualLines != expected.lines || actualWords != expected.words || actualBytes != expected.bytes || actualFile != expected.file {
					t.Errorf("Line %d: expected %d %d %d %s, got %d %d %d %s",
						i+1, expected.lines, expected.words, expected.bytes, expected.file,
						actualLines, actualWords, actualBytes, actualFile)
				}
			}
		})
	}
}

// TestIllegalOption tests the handling of illegal options
func TestIllegalOption(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedErr string
	}{
		{
			name:        "Single Illegal Option",
			args:        []string{"-b"},
			expectedErr: "illegal option -- b",
		},
		{
			name:        "Multiple Options with Illegal",
			args:        []string{"-lwcb"},
			expectedErr: "illegal option -- b",
		},
		{
			name:        "Valid and Invalid Options",
			args:        []string{"-lw", "-x"},
			expectedErr: "illegal option -- x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseArgs(tt.args)
			if err == nil {
				t.Errorf("Expected error, but got nil")
			} else if err.Error() != tt.expectedErr {
				t.Errorf("Expected error: %s, but got: %s", tt.expectedErr, err.Error())
			}
		})
	}
}

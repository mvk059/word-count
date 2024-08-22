# mwc
A Go implementation of the Unix `wc` (word count) command line tool.

## Description
`mwc` is a custom implementation of the Unix `wc` command, written in Go. It provides functionality to count bytes, lines, words, and characters in text files or from standard input. 

## Features

- Count bytes, lines, words, and characters in files
- Read from multiple files
- Read from standard input when no file is specified
- Customizable output based on command-line options
- Handles both ASCII and Unicode text

## Installation

To install `mwc`, make sure you have Go installed on your system, then run:

```
go get github.com/mvk059/mwc
```

## Usage

```
mwc [-lwcm] [file ...]
```

### Options:

- `-l`: Count lines
- `-w`: Count words
- `-c`: Count bytes
- `-m`: Count characters

If no options are specified, `mwc` defaults to counting lines, words, and bytes (equivalent to `-lwc`).

If no filename is provided, `mwc` reads from standard input.

### Examples:

1. Count lines, words, and bytes in a file:
   ```
   mwc file.txt
   ```

2. Count only characters in multiple files:
   ```
   mwc -m file1.txt file2.txt
   ```

3. Count lines from standard input:
   ```
   cat file.txt | mwc -l
   ```

4. Count words and characters in multiple files:
   ```
   mwc -wm file1.txt file2.txt file3.txt
   ```

## Output

The output format is:

```
  <count1> <count2> ... <filename>
```

Where `<count1>`, `<count2>`, etc., are the counts for each specified option, in the order they were provided. The counts are right-aligned in 8-character wide fields.

If multiple files are provided, a total count is displayed at the end.

## Error Handling
- If an invalid option is provided, an error message is displayed, and the program exits.
- If a file cannot be opened or read, an error message is displayed, but the program continues processing other files if any.

## Limitations
- The current implementation may not handle extremely large files efficiently, as it reads the entire file into memory.
- Unicode handling might not be perfect for all edge cases.

## Contributing
Contributions to `mwc` are welcome! Please feel free to submit a Pull Request.

## License
This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments
This project was created as an exercise in Go programming and to understand the inner workings of the `wc` command.
This project was also inspired by and developed as part of the Coding Challenges series. This is my solution to the "Build Your Own wc Tool" challenge.

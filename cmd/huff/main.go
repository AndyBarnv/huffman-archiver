package main

import (
	"flag"
	"fmt"
	"huffman-archiver/internal/archiver"
	"os"
	"path/filepath"
)

func main() {
	compressMode := flag.Bool("c", false, "Comress file")
	decompressMode := flag.Bool("d", false, "Decompress file")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 || (*compressMode == *decompressMode) {
		fmt.Fprintln(os.Stderr, "Usage: huff -c <input> <output> OR huff -d <input> <output>")
		os.Exit(1)
	}

	inputFile, outputFile := args[0], args[1]

	inPath, err := filepath.Abs(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(1)
	outPath, err := filepath.Abs(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(1)

	if inPath == outPath {
		fmt.Fprintf(os.Stderr, "Error: Input and output files cannot be the same (%s)\n", inPath)
		os.Exit(1)
	}

	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	if *compressMode == true {
		err = archiver.Compress(in, out)
	} else {
		err = archiver.Decompress(in, out)
	}

	if err != nil {
		out.Close()
		os.Remove(outputFile)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(1)

	fmt.Println("Done.")
}

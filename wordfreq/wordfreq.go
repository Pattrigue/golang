package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "wordfreq:", err)
		os.Exit(1)
	}
}

func run() error {
	path, n := parseArgs()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	scores, err := countWords(f, n)
	if err != nil {
		return err
	}

	for i := range scores {
		fmt.Printf("%s: %d\n", scores[i].word, scores[i].frequency)
	}

	return nil
}

func parseArgs() (string, int) {
	n := flag.Int("n", 10, "number of words to show")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: wordfreq [-n count] <file>")
		os.Exit(1)
	}

	if *n < 1 {
		fmt.Fprintln(os.Stderr, "-n must be at least 1")
		os.Exit(1)
	}

	path := flag.Arg(0)

	return path, *n
}

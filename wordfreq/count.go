package main

import (
	"bufio"
	"cmp"
	"io"
	"slices"
	"strings"
	"unicode"
)

type WordScore struct {
	word      string
	frequency int
}

func countWords(r io.Reader, n int) ([]WordScore, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	words := map[string]int{}

	for scanner.Scan() {
		// strip leading and trailing characters that aren't 	letters or digits
		word := strings.TrimFunc(scanner.Text(), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		// if now empty, it was a punctuation or symbols
		if word == "" {
			continue
		}

		word = strings.ToLower(word)
		words[word]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	scores := make([]WordScore, 0, len(words))

	for key, value := range words {
		scores = append(scores, WordScore{key, value})
	}

	slices.SortFunc(scores, func(a, b WordScore) int {
		return cmp.Or(
			cmp.Compare(b.frequency, a.frequency),
			cmp.Compare(a.word, b.word),
		)
	})

	return scores[:min(n, len(scores))], nil
}

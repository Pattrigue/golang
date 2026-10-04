package main

import (
	"slices"
	"strings"
	"testing"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		n     int
		want  []WordScore
	}{
		{
			name:  "orders by frequency",
			input: "the cat and the hat and the bat",
			n:     2,
			want:  []WordScore{{"the", 3}, {"and", 2}},
		},
		{
			name:  "ignores case and punctuation",
			input: "Hello, hello! HELLO? world",
			n:     10,
			want:  []WordScore{{"hello", 3}, {"world", 1}},
		},
		{
			name:  "empty input",
			input: "",
			n:     5,
			want:  []WordScore{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countWords(strings.NewReader(tt.input), tt.n)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

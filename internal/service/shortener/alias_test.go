package shortener

import (
	"strings"
	"testing"
)

type tc struct {
	name string
	in   int
	want int
}

func TestRandomLengthPass(t *testing.T) {
	tests := []tc{
		{
			name: "length 6 -> PASS with len(string) = 6",
			in:   6,
			want: 6,
		},
		{
			name: "length 5 -> PASS with len(string) = 5",
			in:   5,
			want: 5,
		},
		{
			name: "length 1 -> PASS with len(string) = 1",
			in:   1,
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Random(tt.in)
			if len(got) != tt.want {
				t.Fatalf("got %s, with len = %d, want len = %d", got, len(got), tt.want)
			}
		})
	}
}

func TestRandomLengthFail(t *testing.T) {
	tests := []tc{
		{
			name: "length 7 -> FAIL with len(string) = 6",
			in:   7,
			want: 6,
		},
		{
			name: "length 6 -> FAIL with len(string) = 5",
			in:   6,
			want: 5,
		},
		{
			name: "length 2 -> FAIL with len(string) = 1",
			in:   2,
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Random(tt.in)
			if len(got) != tt.want {
				t.Fatalf("got %s, with len = %d, want len = %d", got, len(got), tt.want)
			}
		})
	}
}

func TestRandomAllowedChars(t *testing.T) {
	allowed := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for i := 0; i < 100; i++ {
		s := Random(20)
		for _, ch := range s {
			if !strings.ContainsRune(allowed, ch) {
				t.Fatalf("unexpected char: %q in string %q", ch, s)
			}
		}
	}
}

func TestRandomUniqueness(t *testing.T) {
	const n = 1000
	results := make(map[string]struct{})

	for i := 0; i < n; i++ {
		s := Random(10)
		if _, exists := results[s]; exists {
			t.Fatalf("duplicate detected: %s", s)
		}
		results[s] = struct{}{}
	}
}

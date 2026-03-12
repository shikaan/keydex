package cli

import (
	"io"
	"os"
	"testing"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"returns true for y", "y\n", true},
		{"returns true for Y", "Y\n", true},
		{"returns false for n", "n\n", false},
		{"returns false for N", "N\n", false},
		{"returns false for empty input", "\n", false},
		{"returns false for arbitrary input", "foo\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			defer r.Close()

			oldStdin := os.Stdin
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()

			io.WriteString(w, tt.input)
			w.Close()

			if got := Confirm("prompt"); got != tt.want {
				t.Errorf("Confirm() = %v, want %v", got, tt.want)
			}
		})
	}
}

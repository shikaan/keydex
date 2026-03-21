package cli

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestReadSecret_promptGoesToStderr(t *testing.T) {
	rErr, wErr, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = wErr
	defer func() { os.Stderr = oldStderr }()

	rOut, wOut, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = wOut
	defer func() { os.Stdout = oldStdout }()

	ReadSecret("secret-prompt: ")

	wErr.Close()
	wOut.Close()

	stderrContent, _ := io.ReadAll(rErr)
	stdoutContent, _ := io.ReadAll(rOut)
	rErr.Close()
	rOut.Close()

	if !strings.Contains(string(stderrContent), "secret-prompt: ") {
		t.Errorf("expected prompt on stderr, got stderr=%q stdout=%q", stderrContent, stdoutContent)
	}
	if strings.Contains(string(stdoutContent), "secret-prompt: ") {
		t.Errorf("expected prompt NOT on stdout, got %q", stdoutContent)
	}
}

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

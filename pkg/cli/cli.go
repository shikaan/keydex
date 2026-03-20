package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var ReadSecret = func(prompt string) string {
	result := ""
	fmt.Fprint(os.Stderr, prompt)

	for {
		pw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			break
		}
		result = string(pw)
		if result != "" {
			break
		}
	}

	fmt.Fprintln(os.Stderr, "")
	return result
}

var Confirm = func(prompt string) bool {
	fmt.Print(prompt + " [y/N] ")
	out, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}
	out = strings.TrimSpace(out)
	return out == "y" || out == "Y"
}

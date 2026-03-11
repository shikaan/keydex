package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ReadSecret(prompt string) string {
	result := ""
	fmt.Print(prompt)

	for {
		pw, _ := term.ReadPassword(int(syscall.Stdin))
		result = string(pw)
		if result != "" {
			break
		}
	}

	fmt.Println("")
	return result
}

func Confirm(prompt string) bool {
	fmt.Print(prompt + " [y/N] ")
	out, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}
	out = strings.TrimSpace(out)
	return out == "y" || out == "Y"
}

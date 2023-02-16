package credentials

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// Retrieves a locally stored passphrase, if any, otherwise
// prompts the user to insert one
func GetPassphrase(database, passphrase string) string {
	if passphrase != "" {
		return passphrase
	}

	return readFromPrompt(fmt.Sprintf("Passphrase for \"%s\": ", database))
}

func readFromPrompt(promptMessage string) string {
	result := ""
	fmt.Print(promptMessage)

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

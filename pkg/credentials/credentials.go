package credentials

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

const PASSPHRASE = "KPCLI_PASSPHRASE"

// Retrieves a locally stored passphrase, if any, otherwise
// prompts the user to insert one
func GetPassphrase(database string) string {
  passphrase := os.Getenv(PASSPHRASE)

  if passphrase != "" {
    return passphrase
  }

  return readFromPrompt(fmt.Sprintf("Passphrase for \"%s\": ", database))
}

func readFromPrompt(promptMessage string) string {
  result := ""
	fmt.Printf(promptMessage)

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


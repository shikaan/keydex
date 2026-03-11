package credentials

import (
	"fmt"

	"github.com/shikaan/keydex/pkg/cli"
	"github.com/shikaan/keydex/pkg/errors"
)

// Retrieves a locally stored passphrase, if any, otherwise
// prompts the user to insert one
func GetPassphrase(database, passphrase string) string {
	if passphrase != "" {
		return passphrase
	}

	return cli.ReadSecret(fmt.Sprintf("Passphrase for \"%s\": ", database))
}

func MakePassphrase(database string) (string, error) {
	passphrase := cli.ReadSecret(fmt.Sprintf("Create a passphrase for \"%s\": ", database))
	repeated := cli.ReadSecret("Repeat: ")

	if passphrase != repeated {
		return "", errors.MakeError("Passphrase mismatch.", "credentials")
	}

	return passphrase, nil
}

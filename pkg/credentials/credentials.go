package credentials

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"

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

	if passphrase == "" {
		return "", errors.MakeError("Passphrase cannot be empty.", "credentials")
	}

	return passphrase, nil
}

func CreateXMLKeyFileV2(path string) error {
	if _, err := os.Stat(path); err == nil {
		return errors.MakeError("Key file at "+path+" already exists.", "credentials")
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return errors.MakeError("Cannot generate key: "+err.Error(), "credentials")
	}

	hash := sha256.Sum256(key)
	hashPrefix := fmt.Sprintf("%X", hash[:4])
	hexKey := fmt.Sprintf("%X", key)

	data := []byte(
		"<KeyFile>\n" +
			"  <Meta>\n" +
			"    <Version>2.0</Version>\n" +
			"  </Meta>\n" +
			"  <Key>\n" +
			"    <Data Hash=\"" + hashPrefix + "\">" + hexKey + "</Data>\n" +
			"  </Key>\n" +
			"</KeyFile>\n",
	)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return errors.MakeError("Cannot generate key: "+err.Error(), "credentials")
	}

	return nil
}

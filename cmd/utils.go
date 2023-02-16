package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const ENV_DATABASE = "KEYDEX_DATABASE"
const ENV_PASSPHRASE = "KEYDEX_PASSPHRASE"
const ENV_KEY = "KEYDEX_KEY"

// If zero value reference is passed, reads from stdin to get the value
func ReadReferenceFromStdin(maybeReference string) (string, error) {
	if maybeReference != "" {
		return maybeReference, nil
	}

	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

// Cases:
//   DATABASE=lol keydex copy -> gets ref from stdin (blank ref)
//   keydex copy database -> gets ref from stdin (blank ref)
//   DATABASE=lol keydex copy /ref -> OK (db: lol, ref: /ref)
//   keydex copy database /ref -> OK (db: database, ref: /ref)
//   DATABASE=lol keydex copy database /ref -> (db database, /ref)
//   keydex copy /ref -> uses /ref as db and fails

// DATABASE=lol keydex open -> opens list (blank ref)
// keydex copy database -> gets ref from stdin (blank ref)
// DATABASE=lol keydex open /ref -> OK (db: lol, ref: /ref)
// keydex open database /ref -> OK (db: database, ref: /ref)
// DATABASE=lol keydex open database /ref -> (db database, /ref)
// keydex open /ref -> uses /ref as db and then fails
func ReadDatabaseArguments(cmd *cobra.Command, args []string) (database string, reference string, key string) {
	if len(args) == 0 {
		reference = ""
		database = os.Getenv(ENV_DATABASE)
	}

	if len(args) == 1 {
		envDatabase := os.Getenv(ENV_DATABASE)

		if envDatabase != "" {
			database = envDatabase
			reference = args[0]
		} else {
			database = args[0]
			reference = ""
		}
	}

	if len(args) == 2 {
		database = args[0]
		reference = args[1]
	}

	if keyFlag := cmd.Flag("key"); keyFlag != nil {
		key = keyFlag.Value.String()

		if key == "" {
			key = os.Getenv(ENV_KEY)
		}
	}

	return database, reference, key
}

func DatabaseMustBeDefined() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		database, _, _ := ReadDatabaseArguments(cmd, args)

		if database == "" {
			return errors.New("database must be defined")
		}

		return nil
	}
}

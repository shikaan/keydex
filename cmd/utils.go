package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const ENV_DATABASE = "KPCLI_DATABASE"
const ENV_PASSPHRASE = "KPCLI_PASSPHRASE"
const ENV_KEY = "KPCLI_KEY"

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
//   DATABASE=lol kpcli copy -> gets ref from stdin (blank ref)
//   kpcli copy database -> gets ref from stdin (blank ref)
//   DATABASE=lol kpcli copy /ref -> OK (db: lol, ref: /ref)
//   kpcli copy database /ref -> OK (db: database, ref: /ref)
//   DATABASE=lol kpcli copy database /ref -> (db database, /ref)
//   kpcli copy /ref -> uses /ref as db and fails

// DATABASE=lol kpcli open -> opens list (blank ref)
// kpcli copy database -> gets ref from stdin (blank ref)
// DATABASE=lol kpcli open /ref -> OK (db: lol, ref: /ref)
// kpcli open database /ref -> OK (db: database, ref: /ref)
// DATABASE=lol kpcli open database /ref -> (db database, /ref)
// kpcli open /ref -> uses /ref as db and then fails
func ReadDatabaseArguments(cmd *cobra.Command, args []string) (string, string, string) {
	var reference, database string

	if len(args) == 0 {
		reference = ""
		database = os.Getenv(ENV_DATABASE)
	}

	if len(args) == 1 {
		reference = ""
		database = args[0]
	}

	if len(args) == 2 {
		database = args[0]
		reference = args[1]
	}

  key := cmd.Flag("key").Value.String()

  if key == "" {
    key = os.Getenv(ENV_KEY)
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

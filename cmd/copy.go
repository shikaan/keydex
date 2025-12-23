package cmd

import (
	"os"

	"github.com/shikaan/keydex/pkg/clipboard"
	"github.com/shikaan/keydex/pkg/credentials"
	"github.com/shikaan/keydex/pkg/errors"
	"github.com/shikaan/keydex/pkg/info"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/pkg/log"
	"github.com/spf13/cobra"
)

const DEFAULT_FIELD = "password"

var Copy = &cobra.Command{
	Short: "Copies a field of a reference to the clipboard.",
	Long: `Copies a field of a reference to the clipboard.

Reads a 'reference' from the database at 'file' and copies the value of 'field' to the clipboard.

The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the ` + ENV_DATABASE + ` environment variable.
The 'reference' can be passed either as the last argument, or can be read from stdin - to allow piping.
Use the 'list' command to get a list of all the references in the database.

See "Examples" for more details.`,
	Example: `  # Copy the password of the "github" entry in the "coding" group in the "test" database at test.kdbx
  ` + info.NAME + ` copy test.kdbx /test/coding/github

  # Or copy the username instead
  ` + info.NAME + ` copy -f username test.kdbx /test/coding/github

  # Or with stdin
  export ` + ENV_PASSPHRASE + `=${MY_SECRET_PHRASE}
  echo "/test/coding/github" | ` + info.NAME + ` copy test.kdbx

  # Or with stdin and environment variables
  export ` + ENV_PASSPHRASE + `=${MY_SECRET_PHRASE}
  export ` + ENV_DATABASE + `=test.kdbx
  echo "/test/coding/github" | ` + info.NAME + ` copy

  # List entries, browse them with fzf and copy the username to the clipboard
  export ` + ENV_PASSPHRASE + `=${MY_SECRET_PHRASE}
  export ` + ENV_DATABASE + `=test.kdbx

  ` + info.NAME + ` list | fzf | ` + info.NAME + ` copy -f username`,
	Use: "copy [file] [reference]",
	// Initially this command only copied passowrds, hence the aliases.
	// Keeping them around for backwards compatibility.
	Aliases: []string{"cp", "password", "pwd", "copy-password"},
	Args: cobra.MatchAll(
		cobra.MaximumNArgs(2),
		DatabaseMustBeDefined(),
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		database, reference, key := ReadDatabaseArguments(cmd, args)
		field := cmd.Flag("field").Value.String()
		log.Debugf("Using: database %s, reference \"%s\", key \"%s\"", database, reference, key)

		passphrase := credentials.GetPassphrase(database, os.Getenv(ENV_PASSPHRASE))

		return copy(database, key, passphrase, reference, field)
	},
	DisableAutoGenTag: true,
}

func copy(databasePath, keyPath, passphrase, reference, field string) error {
	reference, err := ReadReferenceFromStdin(reference)
	if err != nil {
		return err
	}

	db, err := kdbx.New(databasePath, passphrase, keyPath)
	if err != nil {
		return err
	}

	if entry := db.GetFirstEntryByPath(reference); entry != nil {
		var value string

		if field == DEFAULT_FIELD {
			value = entry.GetPassword()
		} else {
			value = entry.GetContent(field)
		}

		if value == "" {
			return errors.MakeError(`Missing field "`+field+`" in entry "`+reference+`"`, "copy")
		}

		return clipboard.Write(value)
	}

	return errors.MakeError(`Missing entry at "`+reference+`"`, "copy")
}

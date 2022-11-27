package cmd

import (
	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/credentials"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/spf13/cobra"
)

var Copy = &cobra.Command{
	Short: "Copies the password of a reference to the clipboard.",
	Long: `Copies the password of a reference to the clipboard.

Reads a 'reference' from the database at 'file' and copies the password to the clipboard.

The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the KPCLI_DATABASE environment variable.
The 'reference' can be passed either as the last argument, or can be read from stdin - to allow piping.
Use the 'list' command to get a list of all the references in the database.

See "Examples" for more details.`,
	Example: `  # Copy the "github" entry in the "coding" group in the "test" database at test.kdbx
  kpcli copy test.kdbx /test/coding/github

  # Or with stdin
  export KPCLI_PASSPHRASE=${MY_SECRET_PHRASE}
  echo "/test/coding/github" | kpcli copy test.kdbx

  # Or with stdin and environment variables
  export KPCLI_PASSPHRASE=${MY_SECRET_PHRASE}
  export KPCLI_DATABASE=test.kdbx
  echo "/test/coding/github" | kpcli copy

  # List entries, browse them with fzf and copy the result to the clipboard
  export KPCLI_PASSPHRASE=${MY_SECRET_PHRASE}
  export KPCLI_DATABASE=test.kdbx

  kpcli list | fzf | kpcli copy`,
	Use:     "copy [file] [reference]",
	Aliases: []string{"cp", "password", "pwd"},
	Args:    cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		database, reference := ReadDatabaseArguments(args)
		key := cmd.Flag("key").Value.String()
		passphrase := credentials.GetPassphrase(database)

		return copy(database, key, passphrase, reference)
	},
}

func copy(databasePath, keyPath, passphrase, reference string) error {
	reference, err := ReadReferenceFromStdin(reference)
	if err != nil {
		return err
	}

	db, err := kdbx.NewUnlocked(databasePath, passphrase)
	if err != nil {
		return err
	}

	if entry := db.GetFirstEntryByPath(reference); entry != nil {
		return clipboard.Write(entry.GetPassword())
	}

	return errors.MakeError("Missing entry at "+reference, "copy")
}

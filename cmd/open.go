package cmd

import (
	"os"

	"github.com/shikaan/kpcli/pkg/credentials"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/info"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/tui"
	"github.com/spf13/cobra"
)

var Open = &cobra.Command{
	Use:   "open [file] [reference]",
	Short: "Open the entry editor for a reference.",
  Aliases: []string{"edit"},
  Long:  `Open the entry editor for a reference.

Reads a 'reference' from the database at 'file' and opens the editor there. If no reference is passed, it opens a fuzzy search within the editor.

The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the KPCLI_DATABASE environment variable.
The 'reference' can be passed as last argument; if the reference is missing, it opens a fuzzy search.
Use the 'list' command to get a list of all the references in the database.

See "Examples" for more details.`,
  Example: `  # Opens the "github" entry in the "coding" group in the "test" database at test.kdbx
  ` + info.NAME + ` open test.kdbx /test/coding/github
  
  # Open fuzzy search within the test.kdbx database
  ` + info.NAME + ` open test.kdbx

  # Or with environment variables
  export KPCLI_PASSPHRASE=${MY_SECRET_PHRASE}
  export KPCLI_DATABASE=test.kdbx
  ` + info.NAME + ` open

  # List entries, browse them with fzf and edit the result
  export KPCLI_PASSPHRASE=${MY_SECRET_PHRASE}
  export KPCLI_DATABASE=test.kdbx

  ` + info.NAME + ` list | fzf | ` + info.NAME + ` open`,
	RunE: func(cmd *cobra.Command, args []string) error {
    database, reference := ReadDatabaseArguments(args)
    key := cmd.Flag("key").Value.String()
    passphrase := credentials.GetPassphrase(database, os.Getenv(ENV_PASSPHRASE))

		return open(database, key, passphrase, reference)
	},
}

func open(databasePath, keyPath, passphrase, maybeReference string) error {
	reference, err := ReadReferenceFromStdin(maybeReference)
	if err != nil {
		return err
	}

	db, err := kdbx.NewUnlocked(databasePath, passphrase)
	if err != nil {
		return err
	}

  if entry := db.GetFirstEntryByPath(reference); entry != nil { 
    return tui.Run(tui.State{ Entry: entry, Database: db, Reference: reference })
	}

	return errors.MakeError("Missing entry at "+reference, "copy")
}

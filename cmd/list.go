package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/shikaan/keydex/pkg/credentials"
	"github.com/shikaan/keydex/pkg/info"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/pkg/log"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Short: "Lists all the entries in the database",
	Long: `Lists all the entries in the database. 

The list of references - in the form of - /database/group/.../entry will be printed on stadout, allowing for piping.
The 'file' is the the path to the *.kdbx database. It can be passed either as an argument or via the ` + ENV_DATABASE + ` environment variable.
This command can be used in conjuction with tools such like 'fzf' or 'dmenu' to browse the databse and pipe the result to other commands.

See "Examples" for more details.`,
	Use:     "list [file]",
	Aliases: []string{"ls"},
	Args: cobra.MatchAll(
		cobra.MaximumNArgs(1),
		DatabaseMustBeDefined(),
	),
	Example: `  # List all entries of vault.kdbx database
  ` + info.NAME + ` list vault.kdbx

  # List entries, browse them with fzf and copy the result to the clipboard
  export ` + ENV_PASSPHRASE + `=${MY_SECRET_PHRASE}
  export ` + ENV_DATABASE + `=~/vault.kdbx

  ` + info.NAME + ` list | fzf | ` + info.NAME + ` copy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		database, _, key := ReadDatabaseArguments(cmd, args)
    log.Debugf("Using: database %s, key %s", database, key)
		passphrase := credentials.GetPassphrase(database, os.Getenv(ENV_PASSPHRASE))

		return list(database, key, passphrase)
	},
	DisableAutoGenTag: true,
}

func list(database, key, passphrase string) error {
	kdbx, err := kdbx.New(database, passphrase, key)
	if err != nil {
		return err
	}

	entries := kdbx.GetEntryPaths()

	for _, k := range getSortedKeys(entries) {
		fmt.Println(k)
	}

	return nil
}

func getSortedKeys(entries []kdbx.EntryPath) []kdbx.EntryPath {
	less := func(i, j int) bool {
		numberOfSlashesI := len(strings.Split(entries[i], kdbx.PATH_SEPARATOR))
		numberOfSlashesJ := len(strings.Split(entries[j], kdbx.PATH_SEPARATOR))

		// Sort elements in the same group
		if numberOfSlashesI == numberOfSlashesJ {
			return sort.StringsAreSorted([]string{strings.ToLower(entries[i]), strings.ToLower(entries[j])})
		}

		// Show nested entities close to each other
		return numberOfSlashesI > numberOfSlashesJ
	}
	sort.Slice(entries, less)
	return entries
}

package cmd

import (
	"github.com/shikaan/keydex/pkg/info"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   info.NAME,
	Short: "Manage KeePass databases from your terminal.",
	Long: info.NAME + ` is a command line utility to manage KeePass databases. It comes with subcommands for managing the entries and a simple, display-oriented editor inspired by the minimalism of GNU nano.

` + info.NAME + ` can read the following environment variables:

  - ` + ENV_PASSPHRASE + ` 
    When this variable is set, ` + info.NAME + ` will skip the password prompt. It can be replaced by utils such as 'autoexpect'.

  - ` + ENV_DATABASE + `
    Is the path to the *.kbdx database to unlock. Providing 'file' inline overrides this value.

  - ` + ENV_KEY + `
    Is the path to the optional *.key file used to unlock the database. Providing the '--key' flag inline overrides this value.

All the entries are referenced with a path-like reference string shaped like /database/group1/../groupN/entry where 'database' is the database name, 'groupN' is the group name, and 'entry' is the entry title. 

Internally all the entries are referenced by a UUID, however ` + info.NAME + ` will read the first occurrence of a reference in cases of conflicts. Writes are always done via UUID and they are threfore conflict-safe.
    
Some commands make use of the system clipboard, in absence of which ` + info.NAME + ` will fail.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	DisableAutoGenTag: true,
	Version:           info.VERSION,
}

func init() {
	Root.AddCommand(Copy)
	Root.AddCommand(List)
	Root.AddCommand(Open)

	Root.PersistentFlags().StringP("key", "k", "", "path to the key file to unlock the database")
	Copy.Flags().StringP("field", "f", DEFAULT_FIELD, "field whose value will be copied")
	Open.Flags().Bool("read-only", false, "open "+info.NAME+" in read-only mode")
}

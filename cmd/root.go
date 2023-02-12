package cmd

import (
	"github.com/shikaan/kpcli/pkg/info"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   info.NAME,
	Short: "Manage KeePass databases from your terminal.",
	Long: info.NAME + ` is a simple, display-oriented browser and editor for KeePass databases. The user interface is highly inspired by the minimalism of GNU nano: commands are displayed at the bottom of the screen, and context-sensitive help is provided.

Commands are inserted using control-key (^) combinations. For example, "^C" means "Ctrl+C". ` + info.NAME + ` comes with subcommands to read and write entries in the provided database. More information available at "` + info.NAME + ` help [command]". 

To facilitate scripting, this tool comes with the ability of reading the following environment variables:

  - ` + ENV_PASSPHRASE + ` 
    When this variable is set, kpcli will skip the password prompt. It can be replaced by utils such as 'autoexpect'.

  - ` + ENV_DATABASE + `
    Is the path to the *.kbdx database to unlock. Providing 'file' inline overrides this value.

  - ` + ENV_KEY + `
    Is the path to the optional *.key file used to unlock the database. Providing the '--key' flag inline overrides this value.

All the entries are referenced with a path-like reference string shaped like /database/group1/../groupN/entry where 'database' is the database name, 'groupX' is the group name, and 'entry' is the entry title. Internally all the entries are referenced by a UUID, however ` + info.NAME + ` will read the first occurrence of a reference in cases of conflicts. Writes are always done via UUID and they are threfore conflict-safe.
    
Some commands make use of the system clipboard, in absence of which the command will silently fail.

More specific help is available contextually or by typing "` + info.NAME + ` help [command]".`,
	DisableAutoGenTag: true,
}

func init() {
	Root.AddCommand(Copy)
	Root.AddCommand(List)
	Root.AddCommand(Open)

	Root.PersistentFlags().StringP("key", "k", "", "Path to the key file to unlock the database")
}

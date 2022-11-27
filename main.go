package main

import (
  "github.com/spf13/cobra"

	"github.com/shikaan/kpcli/cmd"
)

//go:generate make info


// TODO: document environment variables
func main() {
	var rootCmd = &cobra.Command{
		Use:   "kpcli",
		Short: "Manage KeePass databases from your terminal.",
    Long: `Manage KeePass databases from your terminal.

This utiliy exposes commands to read and write entries in a kdbx database.
To facilitate scripting, this tool comes with the ability of reading the following environment variables:

  - KPCLI_PASSPHRASE 
    When this variable is set, kpcli will skip the password prompt. It can be replaced by utils such as 'autoexpect'.

  - KPCLI_DATABASE
    Is the path to the *.kbdx database to unlock. Providing 'file' inline overrides this value.

All the entries are referenced with a path-like reference string shaped like /database/group1/../groupN/entry where 'database' is the database name, 'groupX' is the group name, and 'entry' is the entry title. Internally all the entries are referenced by a UUID, however kpcli will read the first occurrence of a reference in cases of conflicts. Writes are always done via UUID and they are threfore conflict-safe.
    
Some commands make use of the system clipboard, in absence of which the command will silently fail.
`,
	}

	rootCmd.AddCommand(cmd.Copy)
	rootCmd.AddCommand(cmd.List)
	rootCmd.AddCommand(cmd.Open)

	rootCmd.PersistentFlags().StringP("key", "k", "", "Path to the key file to unlock the database")

	e := rootCmd.Execute()

	if e != nil {
		println(e.Error())
	}
}

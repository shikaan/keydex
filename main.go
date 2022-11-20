package main

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/credentials"
	"github.com/shikaan/kpcli/pkg/kdbx"
	c "github.com/shikaan/kpcli/cmd"
)

//go:generate make info

func main() {	
  var rootCmd = &cobra.Command{
		Use:   "kpcli",
		Short: "Work with KeePass databases from your terminal",
	}
	
  browseCmd.AddCommand(browseCopyCmd)
  
  rootCmd.AddCommand(browseCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
  
  rootCmd.PersistentFlags().StringP("key", "k", "", "Path to the key file to unlock the database")

	e := rootCmd.Execute()

	if e != nil {
		println(e.Error())
	}
}

var listCmd = &cobra.Command{
	Short: "Lists all the entries in the database",
	Long:  "Lists all the entries in the database. Useful to be piped to other tools.",
	Use:   "list DATABASE",
	RunE: func(cmd *cobra.Command, args []string) error {
    database := args[0] 
		key := cmd.Flag("key").Value.String()
    passphrase := credentials.GetPassphrase(database)
		
    return c.List(database, key, passphrase)
	},
}

var copyCmd = &cobra.Command{
	Short: "Copies the password of a reference to the clipboard",
	Long:  "Copies the password for REF in the form /Group/Subgroup/Entry to the clipboard",
	Use:   "copy DATABASE",
	RunE: func(cmd *cobra.Command, args []string) error {
    database := args[0] 
		key := cmd.Flag("key").Value.String()
    passphrase := credentials.GetPassphrase(database)
		
    return c.Copy(database, key, passphrase)
	},
}

var browseCmd = &cobra.Command{
	Short: "Fuzzy search through the entries in DATABASE",
	Use:   "browse [command] DATABASE",
}

var browseCopyCmd = &cobra.Command{
  Use: "copy",
  Short: "Copies entry field",
	RunE: func(cmd *cobra.Command, args []string) error {
		database := args[0]
		key := cmd.Flag("key").Value.String()
    databaseName := filepath.Base(database)
		password := credentials.GetPassphrase(databaseName)

		return c.Browse(database, key, password, func(entry kdbx.Entry) error {
			return clipboard.Write(entry.GetPassword())
    })
	},
}

var editCmd = &cobra.Command{
  Use: "edit DATABASE [REFERENCE]",
  Short: "Edits the entry",
  RunE: func(cmd *cobra.Command, args []string) error {
    database := args[0] 
    maybeRef := ""

    if len(args) == 2 {
      maybeRef = args[1]
    }

		key := cmd.Flag("key").Value.String()
    passphrase := credentials.GetPassphrase(database)

    return c.Edit(database, key, passphrase, maybeRef)
  },
}
package main

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shikaan/kpcli/pkg/credentials"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kpcli",
		Short: "Work with KeePass databases from your terminal",
	}
	

  browseCmd.AddCommand(browseCopyCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(browseCmd)
  
  rootCmd.PersistentFlags().StringP("key", "k", "", "Path to the key file to unlock the database")
  copyCmd.Flags().StringP("field", "f", "password", "field to retrieve")
  browseCopyCmd.Flags().StringP("field", "f", "password", "field to retrieve")

	e := rootCmd.Execute()

	if e != nil {
		println(e.Error())
	}
}

var copyCmd = &cobra.Command{
	Short: "Copies the password of a reference to the clipboard",
	Long:  "Copies the password for REF in the form /Group/Subgroup/Entry to the clipboard",
	Use:   "copy DATABASE",
	RunE: func(cmd *cobra.Command, args []string) error {
    database := args[0] 
		key := cmd.Flag("key").Value.String()
    field := cmd.Flag("field").Value.String()
    passphrase := credentials.GetPassphrase(database)
		
    return Copy(database, key, passphrase, field)
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
		field := cmd.Flag("field").Value.String()
    databaseName := filepath.Base(database)
		password := credentials.GetPassphrase(databaseName)

		return Browse(database, key, password, func(entry kdbx.Entry) error {
			return CopyEntryField(entry, field)
    })
	},
}


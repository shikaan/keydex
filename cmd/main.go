package main

import (
	"fmt"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/credentials"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/spf13/cobra"
)

var key string
var password string

var copyCmd = &cobra.Command{
	Short: "Copies the password of a reference to the clipboard",
	Long:  "Copies the password for REF in the form /Group/Subgroup/Entry to the clipboard",
	Args:  cobra.MinimumNArgs(2),
	Use:   "copy REFERENCE DATABASE",
	RunE: func(cmd *cobra.Command, args []string) error {
		var reference = args[0]
		var database = args[1]

		return Copy(database, "", reference)
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
		password := credentials.GetPassphrase(database)

		return Browse(database, key, password, func(entry kdbx.Entry) error {

			err := clipboard.Write(entry.GetPassword())

			if err != nil {
				return err
			}

			fmt.Printf("Password for \"%s\" copied to clipboard!\n", entry.GetTitle())
			return nil
		})
	},
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kpcli",
		Short: "Work with KeePass databases from your terminal",
	}

	rootCmd.PersistentFlags().StringVar(&key, "key", "", "Path to the key file to unlock the database")

  browseCmd.AddCommand(browseCopyCmd)

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(browseCmd)

	e := rootCmd.Execute()

	if e != nil {
		println(e.Error())
	}
}

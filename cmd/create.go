package cmd

import (
	"os"

	"github.com/shikaan/keydex/pkg/cli"
	"github.com/shikaan/keydex/pkg/credentials"
	"github.com/shikaan/keydex/pkg/errors"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/tui"
	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:     "create [filepath] [name]",
	Short:   "Create an empty KeePass archive.",
	Aliases: []string{"new"},
	Long:    ``,
	Example: ``,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]
		databaseName := args[1]

		if _, err := os.Stat(filepath); err == nil {
			return errors.MakeError("File already exists.", "create")
		}

		passphrase, err := credentials.MakePassphrase(filepath)
		if err != nil {
			return err
		}

		if passphrase == "" {
			return errors.MakeError("Passphrase cannot be empty.", "create")
		}

		file, err := os.Create(filepath)
		if err != nil {
			return errors.MakeError(`Cannot create file: `+err.Error(), "create")
		}

		db, err := kdbx.NewFromFile(file)
		if err != nil {
			return err
		}

		if err = db.SetPasswordAndKey(passphrase, ""); err != nil {
			file.Close()
			return err
		}
		rootGroup := db.NewGroup(databaseName)
		db.Content.Root.Groups = []kdbx.Group{*rootGroup}

		if err = db.Save(); err != nil {
			os.Remove(filepath)
			return err
		}

		if cli.Confirm("Creation successful. Do you want to open the database?") {
			database, err := kdbx.OpenFromPath(filepath, passphrase, "")
			if err != nil {
				return err
			}

			return tui.Run(tui.State{
				Entry:     nil,
				Group:     nil,
				Database:  database,
				Reference: "",
			}, false)
		}

		return nil
	},
	DisableAutoGenTag: true,
}

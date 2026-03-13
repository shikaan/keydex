package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shikaan/keydex/pkg/cli"
	"github.com/shikaan/keydex/pkg/credentials"
	"github.com/shikaan/keydex/pkg/errors"
	"github.com/shikaan/keydex/pkg/info"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/tui"
	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:     "create [file] [name]",
	Short:   "Create an empty KeePass archive.",
	Aliases: []string{"new"},
	Long: `Create an empty KeePass archive.

Creates a new KeePass database at 'file' called 'name'. You will be prompted to set a passphrase for the new database.

See "Examples" for more details.`,
	Example: `  # Create a new database called "vault" at vault.kdbx
  ` + info.NAME + ` create vault.kdbx vault

  # Create a new database at a specific path
  ` + info.NAME + ` create ~/passwords/work.kdbx work`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		name := args[1]

		if _, err := os.Stat(path); err == nil {
			return errors.MakeError("Database file "+path+" already exists.", "create")
		}

		passphrase, err := credentials.MakePassphrase(path)
		if err != nil {
			return err
		}

		keyfilepath := ""
		if cli.Confirm("Do you want to create a keyfile?") {
			filename := strings.Replace(filepath.Base(path), filepath.Ext(path), "-key.xml", 1)
			keyfilepath = filepath.Join(filepath.Dir(path), filename)

			if err = credentials.CreateXMLKeyFileV2(keyfilepath); err != nil {
				return err
			}
		}

		file, err := os.Create(path)
		if err != nil {
			if keyfilepath != "" {
				_ = os.Remove(keyfilepath)
			}
			return errors.MakeError(`Cannot create file: `+err.Error(), "create")
		}

		defer file.Close()

		db, err := kdbx.NewFromFile(file)
		if err != nil {
			if keyfilepath != "" {
				_ = os.Remove(keyfilepath)
			}
			os.Remove(path)
			return err
		}

		if err = db.SetPasswordAndKey(passphrase, keyfilepath); err != nil {
			if keyfilepath != "" {
				_ = os.Remove(keyfilepath)
			}
			os.Remove(path)
			return err
		}
		rootGroup := db.NewGroup(name)
		db.Content.Root.Groups = []kdbx.Group{*rootGroup}

		db.Database.Content.Meta.DatabaseName = name

		if err = db.SaveAndUnlockEntries(); err != nil {
			if keyfilepath != "" {
				_ = os.Remove(keyfilepath)
			}
			os.Remove(path)
			return err
		}

		if cli.Confirm("Creation successful. Do you want to open the database?") {
			return tui.Run(tui.State{
				Entry:     nil,
				Group:     nil,
				Database:  db,
				Reference: "",
			}, false)
		}

		return nil
	},
	DisableAutoGenTag: true,
}

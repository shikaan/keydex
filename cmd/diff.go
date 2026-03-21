package cmd

import (
	"fmt"
	"os"

	"github.com/shikaan/keydex/pkg/credentials"
	"github.com/shikaan/keydex/pkg/info"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/pkg/log"
	"github.com/spf13/cobra"
)

var Diff = &cobra.Command{
	Short: "Compares two KeePass archives",
	Long: `Compares two KeePass archives and outputs which entries were added, removed,
or modified. Output follows the unified diff format so it can be piped into other tools.

The 'file-a' and 'file-b' arguments are paths to the *.kdbx archives to compare.
Passphrases can be provided via environment variables to avoid interactive prompts.`,
	Use: "diff [file-a] [file-b]",
	Args: cobra.ExactArgs(2),
	Example: `  # Compare two archives
  ` + info.NAME + ` diff old.kdbx new.kdbx

  # Or with environment variables
  export ` + ENV_PASSPHRASE_A + `=${PASSPHRASE_A}
  export ` + ENV_PASSPHRASE_B + `=${PASSPHRASE_B}
  ` + info.NAME + ` diff old.kdbx new.kdbx

  # With key files
  ` + info.NAME + ` diff --key-a old.key --key-b new.key old.kdbx new.kdbx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileA, fileB := args[0], args[1]

		keyA, _ := cmd.Flags().GetString("key-a")
		keyB, _ := cmd.Flags().GetString("key-b")

		log.Infof("Using: file-a: %s, key-a: %s, file-b: %s, key-b: %s",
			fileA, orDefault(keyA), fileB, orDefault(keyB))

		passphraseA := credentials.GetPassphrase(fileA, os.Getenv(ENV_PASSPHRASE_A))
		passphraseB := credentials.GetPassphrase(fileB, os.Getenv(ENV_PASSPHRASE_B))

		return diff(fileA, fileB, keyA, keyB, passphraseA, passphraseB)
	},
	DisableAutoGenTag: true,
}

func diff(fileA, fileB, keyA, keyB, passphraseA, passphraseB string) error {
	dbA, err := kdbx.OpenFromPath(fileA, passphraseA, keyA)
	if err != nil {
		return err
	}

	dbB, err := kdbx.OpenFromPath(fileB, passphraseB, keyB)
	if err != nil {
		return err
	}

	diffs := kdbx.DiffDatabases(dbA, dbB)
	fmt.Print(kdbx.FormatDiff(fileA, fileB, diffs))

	return nil
}

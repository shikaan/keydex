package cmd

import (
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/pkg/tui"
)

func Edit(databasePath, keyPath, passphrase, maybeReference string) error {
	reference, err := readReferenceFromStdin(maybeReference)
	if err != nil {
		return err
	}

	db, err := kdbx.NewUnlocked(databasePath, passphrase)
	if err != nil {
		return err
	}

  if entry := db.GetEntry(reference); entry != nil { 
    return tui.Run(tui.State{ Entry: entry, Database: db, Reference: reference })
	}

	return errors.MakeError("Missing entry at "+reference, "copy")
}

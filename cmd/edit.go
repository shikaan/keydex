package main

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

  entries := db.GetEntries()

	if entry, ok := entries[reference]; ok {
		return tui.RunEditView(*entry)
	}

	return errors.MakeError("Missing entry at "+reference, "copy")
}

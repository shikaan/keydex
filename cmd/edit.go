package main

import (
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/tui"
	"github.com/shikaan/kpcli/pkg/kdbx"
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

  if entry, ok := db.Entries[reference]; ok {
    return tui.OpenEntryEditor(entry)
  }

  return errors.MakeError("Missing entry at " + reference, "copy") 
}

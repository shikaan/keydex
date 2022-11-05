package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

// Reads reference from stdin and attempts to copy password
// of referenced entry to the clipboard
func Copy(databasePath, keyPath, passphrase string) error { 
  reference, err := readReferenceFromStdin("")
  if err != nil {
    return err
  }
  
  db, err := kdbx.NewUnlocked(databasePath, passphrase)
  if err != nil {
    return err
  }

  if entry, ok := db.Entries[reference]; ok {
    return clipboard.Write(entry.Password)
  }

  return errors.MakeError("Missing entry at " + reference, "copy") 
}

func readReferenceFromStdin(maybeReference string) (string, error) {
  if maybeReference != "" {
    return maybeReference, nil
  } 

  reader := bufio.NewReader(os.Stdin)
  str, err := reader.ReadString('\n')

  if err != nil {
    return "", err
  }

  return strings.TrimSpace(str), nil
}

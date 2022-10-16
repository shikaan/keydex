package main

import (
	"bufio"
	"os"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

// Reads reference from stdin and attempts to copy password
// of referenced entry to the clipboard
func Copy(databasePath, keyPath, password string) error { 
  reference := readReferenceFromStdin()

  kdbx, err := kdbx.New(databasePath)
	if err != nil {
		return err
	}

	err = kdbx.Unlock(password)
	if err != nil {
		return err
	}

  if entry, ok := kdbx.Entries[reference]; ok {
    clipboard.Write(entry.GetPassword())
    return nil
  }

  return errors.MakeError("Unable to find entry at " + reference, "copy") 
}

func readReferenceFromStdin() string {
  value := ""
  scanner := bufio.NewScanner(os.Stdin)
  
  for scanner.Scan() {
    value = value + scanner.Text()
  }

  return value
}

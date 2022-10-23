package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

// Reads reference from stdin and attempts to copy password
// of referenced entry to the clipboard
func Copy(databasePath, keyPath, passphrase, field string) error { 
  reference := readReferenceFromStdin()

  println(reference)

  db, err := kdbx.NewUnlocked(databasePath, passphrase)
  if err != nil {
    return err
  }

  if entry, ok := db.Entries[reference]; ok {
    return CopyEntryField(entry, field)
  }

  return errors.MakeError("Missing entry at " + reference, "copy") 
}

func CopyEntryField(entry kdbx.Entry, field string) error {
  if content, ok := entry.Fields[field]; ok {
    return clipboard.Write(content)
  }
 
  fields := make([]string, 0, len(entry.Fields))
  for k := range entry.Fields {
    fields = append(fields, k)
  }

  msg := fmt.Sprintf("Missing field %s on %s. Allowed fields: %s.", field, entry.Name, strings.Join(fields, ","))
  return errors.MakeError(msg, "copy") 
}

func readReferenceFromStdin() string {
  val := ""
  scanner := bufio.NewScanner(os.Stdin)
  
  for {
    scanner.Scan()
    val = val + scanner.Text()
    break
  }

  return val
}

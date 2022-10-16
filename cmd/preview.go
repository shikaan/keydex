package main

import (
	"fmt"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

func Preview(databasePath, keyPath, password string) error {
	kdbx, err := kdbx.NewUnlocked(databasePath, password)
  if err != nil {
    return err
  }

  keys := make([]string, 0, len(kdbx.Entries)) 
  for k := range kdbx.Entries {
    keys = append(keys, k)
  }

  idx, _ := fuzzyfinder.Find(
    keys,
    func(i int) string {
      return keys[i]
    },
    fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
      if i == -1 {
        return ""
      }

      if entry, ok := kdbx.Entries[keys[i]]; ok {
        return renderEntry(entry)
      }

      return ""
    }),
  )

  reference := keys[idx]
  if entry, ok := kdbx.Entries[reference]; ok {
    clipboard.Write(entry.GetPassword())
    return nil
  }

  return errors.MakeError("Unable to find entry at " + reference, "browse") 
}

func renderEntry(e kdbx.Entry) string {
  s := &strings.Builder{}

  s.WriteString(e.GetTitle())
  s.WriteString("\n")

  for i, v := range e.Values {
    if v.Key == "Title" || v.Value.Content == "" {
      continue
    }
    
    s.WriteString("\n")
    if i == e.GetPasswordIndex() {
      s.WriteString("Password: ***") 
    } else {
      s.WriteString(fmt.Sprintf("%s: %s", v.Key, v.Value.Content))
    }
  }

  return s.String()
}

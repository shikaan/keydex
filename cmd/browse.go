package main

import (
	"fmt"
	"sort"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

func Browse(database, key, passphrase string, callback func (entry kdbx.Entry) error) error {
	kdbx, err := kdbx.NewUnlocked(database, passphrase)
  if err != nil {
    return err
  }

  entries := kdbx.GetEntries()
  keys := make([]string, 0, len(entries)) 
  for k := range entries {
    keys = append(keys, k)
  }
  sort.Slice(keys, func(i, j int) bool {
    return keys[i] > keys[j]
  })

  idx, err := fuzzyfinder.Find(
    keys,
    func(i int) string {
      return keys[i]
    },
    fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
      if i == -1 {
        return ""
      }

      if entry, ok := entries[keys[i]]; ok {
        return previewEntry(*entry)
      }

      return ""
    }),
  )

  if err != nil {
    return err
  }

  reference := keys[idx]
  if entry, ok := entries[reference]; ok {
    return callback(*entry)
  }

  return errors.MakeError("Unable to find entry at " + reference, "browse") 
}

func previewEntry(e kdbx.Entry) string {
  s := &strings.Builder{}

  s.WriteString(e.GetTitle())
  s.WriteString("\n")

  for i, v := range e.Values {
    if v.Key == "Title" || v.Value.Content == "" {
      continue
    }
    
    isPassword := i == e.GetPasswordIndex()

    s.WriteString("\n")
    if isPassword {
      s.WriteString("Password: ***") 
    } else {
      s.WriteString(fmt.Sprintf("%s: %s", v.Key, v.Value.Content))
    }
  }

  return s.String()
}

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

  keys := make([]string, 0, len(kdbx.Entries)) 
  for k := range kdbx.Entries {
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

      if entry, ok := kdbx.Entries[keys[i]]; ok {
        return previewEntry(entry)
      }

      return ""
    }),
  )

  if err != nil {
    return err
  }

  reference := keys[idx]
  if entry, ok := kdbx.Entries[reference]; ok {
    return callback(entry)
  }

  return errors.MakeError("Unable to find entry at " + reference, "browse") 
}

func previewEntry(e kdbx.Entry) string {
  s := &strings.Builder{}

  s.WriteString(e.Name)
  s.WriteString("\n")

  for _, f := range e.Fields {
    k := f[0]
    v := f[1]
    if k == "Title" || v == "" {
      continue
    }
    
    s.WriteString("\n")
    if k == kdbx.PASSWORD_KEY {
      s.WriteString("Password: ***") 
    } else {
      s.WriteString(fmt.Sprintf("%s: %s", k, v))
    }
  }

  return s.String()
}

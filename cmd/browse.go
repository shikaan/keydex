package cmd

import (
	"fmt"
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

  keys := kdbx.GetEntryPaths()

  idx, err := fuzzyfinder.Find(
    keys,
    func(i int) string {
      return keys[i]
    },
    fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
      if i == -1 {
        return ""
      }

      if entry := kdbx.GetFirstEntryByPath(keys[i]); entry != nil {
        return previewEntry(*entry)
      }

      return ""
    }),
  )

  if err != nil {
    return err
  }

  reference := keys[idx]
 
  if entry := kdbx.GetFirstEntryByPath(reference); entry != nil {
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

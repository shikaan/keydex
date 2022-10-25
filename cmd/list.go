package main

import (
	"fmt"

	"github.com/shikaan/kpcli/pkg/kdbx"
)

func List(database, key, passphrase string) error {
	kdbx, err := kdbx.NewUnlocked(database, passphrase)
  if err != nil {
    return err
  }

  for k := range kdbx.Entries {
    fmt.Println(k)
  }

  return nil
}

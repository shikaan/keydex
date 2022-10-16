package main

import (
	"fmt"

	"github.com/shikaan/kpcli/pkg/kdbx"
)

func List(databasePath, keyPath, password string) error {
	kdbx, err := kdbx.New(databasePath)
	if err != nil {
		return err
	}

	err = kdbx.Unlock(password)
	if err != nil {
		return err
	}

  for k := range kdbx.Entries {
    fmt.Println(k)
  }

  return nil
}


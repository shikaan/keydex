package main

import (
	"fmt"

	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/pkg/logger"
)

func List(database, keyPath, password string, logger *logger.Logger) error {
	kdbx, err := kdbx.New(database)
	if err != nil {
		return err
	}

	err = kdbx.Unlock(password)
	if err != nil {
		return err
	}

  paths, _ := kdbx.ListPaths()

  for _, p := range paths {
    fmt.Println(p)
  }

  return nil
}


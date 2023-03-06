package main

import (
	"fmt"
	"os"

	"github.com/shikaan/keydex/cmd"
	"github.com/shikaan/keydex/pkg/errors"
	"github.com/shikaan/keydex/pkg/log"
)

//go:generate make docs

func main() {
	defer func() {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case error:
				log.Error("Unexpected error", e)
			default:
				log.Error(fmt.Sprintf("Unexpected error: %v", e), nil)
			}
			println(errors.MakeError("An unexpected error occurred. Check logs for details.", "open").Error())

			os.Exit(1)
		}
	}()

	if err := cmd.Root.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/shikaan/kpcli/cmd"
)

//go:generate make docs

func main() {
	if err := cmd.Root.Execute(); err != nil {
		os.Exit(1)
	}
}

//go:build exclude

package main

import (
	"log"

	"github.com/shikaan/keydex/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.Root, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}

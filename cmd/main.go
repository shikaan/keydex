package main

import (
	"github.com/shikaan/kpcli/pkg/open"
	"github.com/spf13/cobra"
)

func main() {

	var keyPath string

	open := &cobra.Command{
		Use:   "open [archive path]",
		Short: "Open specified archive",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			open.Open(args[0], keyPath)
		},
	}

	open.PersistentFlags().StringVar(&keyPath, "key", "k", "path to the key file")

	root := &cobra.Command{Use: "kpcli"}
	root.AddCommand(open)

	root.Execute()
}

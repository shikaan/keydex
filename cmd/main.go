package main

import (
	"github.com/shikaan/kpcli/pages/home"
	"github.com/shikaan/kpcli/pkg/logger"
	"github.com/spf13/cobra"
)

func main() {
	l := logger.NewFileLogger(logger.Debug, "kpcli.log")
	defer l.CleanUp()

	var keyPath string

	open := &cobra.Command{
		Use:   "open [archive path]",
		Short: "Open specified archive",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			home.Open(args[0], keyPath)
		},
	}

	open.PersistentFlags().StringVar(&keyPath, "key", "k", "path to the key file")

	root := &cobra.Command{Use: "kpcli"}
	root.AddCommand(open)

	root.Execute()
}

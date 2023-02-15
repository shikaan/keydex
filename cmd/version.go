package cmd

import (
	"fmt"

	"github.com/shikaan/kpcli/pkg/info"
	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:     "version",
	Short:   fmt.Sprintf("Print the version number of %s.", info.NAME),
	Aliases: []string{"version"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(info.VERSION)
		return nil
	},
	DisableAutoGenTag: true,
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is a version number.
var version = "0.4.0"

func init() {
	RootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version)
		},
	}

	return cmd
}

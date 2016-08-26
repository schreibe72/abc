package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Githash string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "shows the version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Githash: %s\n", Githash)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

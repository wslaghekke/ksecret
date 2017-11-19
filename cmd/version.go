package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Rev = "Unknown"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ksecret",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Rev)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

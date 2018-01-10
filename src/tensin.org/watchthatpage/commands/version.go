package commands

import (
	"fmt"

	"tensin.org/watchthatpage/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of watchthatpage",
	Long:  `Print the version number of watchthatpage`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(core.Version)
	},
}

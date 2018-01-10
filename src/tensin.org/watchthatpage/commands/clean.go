package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tensin.org/watchthatpage/core"
)

func init() {
	RootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cached content",
	Long:  `Clean cached content`,
	Run: func(cmd *cobra.Command, args []string) {
		const path = "cache/"
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Now cleaning content from [" + path + "]")
			err = core.RemoveAllContentFromDirectory(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Path [" + path + "] doesn't exist, nothing to do")
		}
	},
}

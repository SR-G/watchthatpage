package commands

import (
	"github.com/spf13/cobra"
)

// ConfigurationFileName is the configuration file name that will be used
var ConfigurationFileName string

// RootCmd is the main command = the program itself
var RootCmd = &cobra.Command{
	Use:   "watchthatpage",
	Short: "WatchThatPage is a command line program used to trigger notifications when some HTML page contents is modified",
	Long:  `WatchThatPage is a command line program used to trigger notifications when some HTML page contents is modified`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&ConfigurationFileName, "configuration", "", "watchthatpage.json", "Configuration file name. Default is binary name + .json (e.g. 'watchthatpage.json'), in the same folder than the binary itself")
}

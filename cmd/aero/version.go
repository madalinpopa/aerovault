package main

import (
	"github.com/madalinpopa/aerovault"
	"github.com/spf13/cobra"
)

// versionCmd is a cobra command that prints the version number of the CLI when executed.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(aerovault.GetVersion())
	},
}

// init adds the versionCmd subcommand to the rootCmd, enabling the CLI to print version information.
func init() {
	rootCmd.AddCommand(versionCmd)
}

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// rootCmd is the base command for the CLI application. It prints help information by default when no subcommands are provided.
var rootCmd = &cobra.Command{
	Use:   "Usage: aero <command> <args>",
	Short: "Simple CLI to backup and restore Docker volumes",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute runs the root command and handles any errors that occur during execution.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

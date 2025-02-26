/*
Copyright © 2025 EspressoCake
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bloodhound_community",
	Short: "Generation utility for creating Bloodhound Docker instances and fetching metadata",
	Long:  `Generation utility for creating Bloodhound Docker instances and fetching metadata`,
}

func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}

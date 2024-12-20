package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "CLI for managing the backend",
}

func Execute() error {
	return rootCmd.Execute()
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backend-go",
	Short: "Backend application with API and Consumer components",
	Long:  `A complete backend application with API REST and message consumer`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

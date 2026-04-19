package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var envFilePath string

var rootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "A lightweight CLI for managing .env files across environments and profiles",
	Long: `envoy helps you read, write, and manage .env files
across multiple environments and named profiles.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&envFilePath, "file", "f", ".env",
		"path to the .env file",
	)

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
}

// shared helper used by sub-commands
func envPath() string { return envFilePath }

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/envfile"
)

var requiredKeys string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a .env file against required keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		var required []string
		if requiredKeys != "" {
			for _, k := range strings.Split(requiredKeys, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					required = append(required, k)
				}
			}
		}

		result := envfile.Validate(entries, required)
		if result.OK() {
			fmt.Println("✓ validation passed")
			return nil
		}

		fmt.Fprintln(os.Stderr, "✗ validation failed:")
		fmt.Fprintln(os.Stderr, result.String())
		os.Exit(1)
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&requiredKeys, "require", "r", "", "comma-separated list of required keys")
	rootCmd.AddCommand(validateCmd)
}

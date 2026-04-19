package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var diffCmd = &cobra.Command{
	Use:   "diff [base] [other]",
	Short: "Show differences between two .env files or profiles",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseEntries, err := loadSource(args[0])
		if err != nil {
			return fmt.Errorf("loading base: %w", err)
		}
		otherEntries, err := loadSource(args[1])
		if err != nil {
			return fmt.Errorf("loading other: %w", err)
		}

		diffs := envfile.Diff(baseEntries, otherEntries)
		output := envfile.FormatDiff(diffs)
		if output == "" {
			fmt.Println("No differences found.")
		} else {
			fmt.Print(output)
		}
		return nil
	},
}

// loadSource loads entries from a file path or a named profile.
func loadSource(src string) ([]envfile.Entry, error) {
	// Try as file first
	if _, err := os.Stat(src); err == nil {
		return envfile.Parse(src)
	}
	// Fall back to profile
	entries, err := envfile.LoadProfile(src)
	if err != nil {
		return nil, fmt.Errorf("%q is neither a valid file nor a profile", src)
	}
	return entries, nil
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

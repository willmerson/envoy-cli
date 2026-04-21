package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envoy-cli/envoy/internal/envfile"
)

// compareCmd compares two .env files and reports keys that are only in one,
// missing from one, or have differing values between them.
var compareCmd = &cobra.Command{
	Use:   "compare <fileA> <fileB>",
	Short: "Compare two .env files side by side",
	Long: `Compare two .env files and display a structured report showing:
  - Keys only present in file A
  - Keys only present in file B
  - Keys present in both but with different values
  - Keys that are identical in both files`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileA := args[0]
		fileB := args[1]

		entriesA, err := loadEnvFile(fileA)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fileA, err)
		}

		entriesB, err := loadEnvFile(fileB)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fileB, err)
		}

		result := envfile.Compare(entriesA, entriesB)

		showIdentical, _ := cmd.Flags().GetBool("show-identical")
		output := envfile.FormatCompare(result, fileA, fileB, showIdentical)
		fmt.Print(output)

		// Exit with code 1 if there are any differences, useful for CI pipelines
		if exitOnDiff, _ := cmd.Flags().GetBool("exit-on-diff"); exitOnDiff {
			if len(result.OnlyInA) > 0 || len(result.OnlyInB) > 0 || len(result.Different) > 0 {
				os.Exit(1)
			}
		}

		return nil
	},
}

// loadEnvFile reads and parses an .env file from the given path.
func loadEnvFile(path string) ([]envfile.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return envfile.Parse(f)
}

func init() {
	compareCmd.Flags().Bool("show-identical", false, "Include keys that are identical in both files")
	compareCmd.Flags().Bool("exit-on-diff", false, "Exit with code 1 if any differences are found (useful for CI)")
	rootCmd.AddCommand(compareCmd)
}

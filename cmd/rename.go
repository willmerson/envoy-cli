package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/envfile"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old-key> <new-key>",
	Short: "Rename a key in the .env file",
	Long:  `Rename an existing key to a new name while preserving its value and position.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldKey := args[0]
		newKey := args[1]

		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		updated, result, err := envfile.RenameKey(entries, oldKey, newKey)
		if err != nil {
			return fmt.Errorf("rename failed: %w", err)
		}

		if err := envfile.Write(envPath, updated); err != nil {
			return fmt.Errorf("failed to write env file: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Renamed %q → %q\n", result.OldKey, result.NewKey)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

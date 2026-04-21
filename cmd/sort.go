package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	sortDescending bool
	sortInPlace    bool
	sortShowGroups bool
)

// sortCmd sorts the entries in a .env file alphabetically by key.
var sortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "Sort entries in a .env file alphabetically by key",
	Long: `Sort all key-value entries in a .env file alphabetically.

By default, entries are sorted in ascending order. Use --desc for descending.
Use --in-place to overwrite the source file, or output is printed to stdout.
Use --groups to display a summary of prefix-grouped keys after sorting.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine the source file
		src := envPath
		if len(args) > 0 {
			src = args[0]
		}

		// Parse the env file
		entries, err := envfile.Parse(src)
		if err != nil {
			return fmt.Errorf("failed to parse env file %q: %w", src, err)
		}

		// Determine sort order
		order := envfile.SortAscending
		if sortDescending {
			order = envfile.SortDescending
		}

		// Sort the entries
		result, movedCount := envfile.Sort(entries, order)

		// Display group summary if requested
		if sortShowGroups {
			groups := envfile.GroupByPrefix(result)
			fmt.Fprintln(cmd.OutOrStdout(), "Key groups by prefix:")
			for prefix, keys := range groups {
				fmt.Fprintf(cmd.OutOrStdout(), "  [%s]: %d key(s)\n", prefix, len(keys))
			}
			fmt.Fprintln(cmd.OutOrStdout())
		}

		// Write sorted entries
		if sortInPlace {
			if err := envfile.Write(src, result); err != nil {
				return fmt.Errorf("failed to write sorted file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Sorted %q (%d entries reordered)\n", src, movedCount)
		} else {
			if err := envfile.WriteTo(os.Stdout, result); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			if movedCount > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "# %d entries reordered\n", movedCount)
			}
		}

		return nil
	},
}

func init() {
	sortCmd.Flags().BoolVar(&sortDescending, "desc", false, "Sort keys in descending order")
	sortCmd.Flags().BoolVarP(&sortInPlace, "in-place", "i", false, "Overwrite the source file with sorted output")
	sortCmd.Flags().BoolVar(&sortShowGroups, "groups", false, "Show a summary of keys grouped by prefix")
	rootCmd.AddCommand(sortCmd)
}

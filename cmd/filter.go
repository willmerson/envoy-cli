package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	filterPrefix    string
	filterSuffix    string
	filterContains  string
	filterKeysOnly  bool
	filterEmptyOnly bool
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter entries from a .env file by prefix, suffix, or pattern",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		opts := envfile.FilterOptions{
			Prefix:    filterPrefix,
			Suffix:    filterSuffix,
			Contains:  filterContains,
			KeysOnly:  filterKeysOnly,
			EmptyOnly: filterEmptyOnly,
		}

		result := envfile.Filter(entries, opts)

		if len(result.Matched) == 0 {
			fmt.Fprintf(os.Stderr, "No entries matched the filter criteria (checked %d entries).\n", result.Total)
			return nil
		}

		for _, e := range result.Matched {
			if filterKeysOnly {
				fmt.Println(e.Key)
			} else {
				fmt.Printf("%s=%s\n", e.Key, e.Value)
			}
		}

		fmt.Fprintf(os.Stderr, "\nMatched %d of %d entries.\n", len(result.Matched), result.Total)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().StringVar(&filterPrefix, "prefix", "", "Filter keys starting with this prefix")
	filterCmd.Flags().StringVar(&filterSuffix, "suffix", "", "Filter keys ending with this suffix")
	filterCmd.Flags().StringVar(&filterContains, "contains", "", "Filter entries containing this string")
	filterCmd.Flags().BoolVar(&filterKeysOnly, "keys-only", false, "Match and print keys only (not values)")
	filterCmd.Flags().BoolVar(&filterEmptyOnly, "empty-only", false, "Only return entries with empty values")
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

func init() {
	var separator string
	var depth int
	var ungrouped string
	var outputFormat string

	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Group .env entries by key prefix",
		Long:  `Partition entries by a common key prefix (e.g. DB_, APP_) and display them grouped.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("failed to parse env file: %w", err)
			}

			opts := envfile.GroupOptions{
				Separator: separator,
				Depth:     depth,
				Ungrouped: ungrouped,
			}

			result := envfile.Group(entries, opts)

			switch strings.ToLower(outputFormat) {
			case "keys":
				for _, name := range result.Ordered {
					fmt.Printf("%s (%d)\n", name, len(result.Groups[name]))
				}
			case "summary", "":
				fmt.Print(envfile.FormatGroupSummary(result))
			default:
				return fmt.Errorf("unknown output format %q (use: summary, keys)", outputFormat)
			}

			return nil
		},
	}

	groupCmd.Flags().StringVarP(&separator, "separator", "s", "_", "Key segment separator")
	groupCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Number of prefix segments to use for grouping")
	groupCmd.Flags().StringVarP(&ungrouped, "ungrouped", "u", "other", "Label for entries with no matching prefix")
	groupCmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format: summary, keys")

	rootCmd.AddCommand(groupCmd)

	_ = os.Stderr
}

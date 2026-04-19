package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	mergeStrategy string
	mergeOutput   string
)

var mergeCmd = &cobra.Command{
	Use:   "merge <base-profile> <override-profile>",
	Short: "Merge two profiles into one",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := envfile.LoadProfile(args[0])
		if err != nil {
			return fmt.Errorf("loading base profile: %w", err)
		}
		override, err := envfile.LoadProfile(args[1])
		if err != nil {
			return fmt.Errorf("loading override profile: %w", err)
		}

		strategy := envfile.StrategyOurs
		if mergeStrategy == "theirs" {
			strategy = envfile.StrategyTheirs
		}

		result := envfile.Merge(base, override, strategy)
		fmt.Fprintln(os.Stderr, "Merge summary:", envfile.MergeSummary(result))

		if mergeOutput != "" {
			return envfile.SaveProfile(mergeOutput, result.Entries)
		}
		return envfile.WriteTo(os.Stdout, result.Entries)
	},
}

func init() {
	mergeCmd.Flags().StringVar(&mergeStrategy, "strategy", "ours", "Conflict strategy: ours|theirs")
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "Save result to a profile instead of stdout")
	rootCmd.AddCommand(mergeCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	rollbackDryRun  bool
	rollbackVerbose bool
)

func init() {
	rollbackCmd := &cobra.Command{
		Use:   "rollback <snapshot-name>",
		Short: "Restore a .env file to a previously saved snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapName := args[0]

			currentEntries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}

			snapEntries, err := envfile.LoadSnapshot(snapName)
			if err != nil {
				return fmt.Errorf("loading snapshot %q: %w", snapName, err)
			}

			if rollbackDryRun {
				plan := envfile.PlanRollback(currentEntries, snapEntries)
				fmt.Println("Dry-run rollback plan:")
				printKeys("  restore", plan.ToRestore)
				printKeys("  add    ", plan.ToAdd)
				printKeys("  remove ", plan.ToRemove)
				return nil
			}

			restored, result := envfile.Rollback(currentEntries, snapEntries)
			result.SnapshotName = snapName

			if err := envfile.Write(envPath, restored); err != nil {
				return fmt.Errorf("writing env file: %w", err)
			}

			fmt.Println(envfile.FormatRollbackResult(result))
			if rollbackVerbose {
				plan := envfile.PlanRollback(currentEntries, snapEntries)
				printKeys("  restored", plan.ToRestore)
				printKeys("  added   ", plan.ToAdd)
				printKeys("  removed ", plan.ToRemove)
			}
			return nil
		},
	}

	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "Preview changes without applying them")
	rollbackCmd.Flags().BoolVarP(&rollbackVerbose, "verbose", "v", false, "Show detailed key changes")
	rootCmd.AddCommand(rollbackCmd)
}

func printKeys(label string, keys []string) {
	if len(keys) == 0 {
		return
	}
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%s: %s\n", label, k)
	}
}

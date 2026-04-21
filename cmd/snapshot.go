package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage env file snapshots",
	}

	// snapshot save
	saveCmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save a snapshot of the current .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("failed to parse env file: %w", err)
			}
			if err := envfile.SaveSnapshot(name, envPath, entries); err != nil {
				return fmt.Errorf("failed to save snapshot: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot %q saved (%d entries)\n", name, len(entries))
			return nil
		},
	}

	// snapshot restore
	restoreCmd := &cobra.Command{
		Use:   "restore <name>",
		Short: "Restore a snapshot to the current .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rec, err := envfile.LoadSnapshot(args[0])
			if err != nil {
				return err
			}
			if err := envfile.Write(envPath, rec.Entries); err != nil {
				return fmt.Errorf("failed to write env file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Restored snapshot %q (%d entries) to %s\n",
				rec.Meta.Name, len(rec.Entries), envPath)
			return nil
		},
	}

	// snapshot list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			metas, err := envfile.ListSnapshots()
			if err != nil {
				return err
			}
			if len(metas) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No snapshots found.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tFILE\tENTRIES\tCREATED")
			for _, m := range metas {
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
					m.Name, m.File, m.EntryCount,
					m.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			return w.Flush()
		},
	}

	snapshotCmd.AddCommand(saveCmd, restoreCmd, listCmd)
	rootCmd.AddCommand(snapshotCmd)
}

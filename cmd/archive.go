package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/envfile"
)

func init() {
	archiveLabel := ""
	archiveID := ""

	archiveCmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive and restore env file versions",
	}

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Save current env file as a named archive",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			base, _ := os.Getwd()
			id, err := envfile.SaveArchive(base, archiveLabel, entries)
			if err != nil {
				return err
			}
			fmt.Printf("Archive saved: %s (label: %q)\n", id, archiveLabel)
			return nil
		},
	}
	saveCmd.Flags().StringVarP(&archiveLabel, "label", "l", "", "Label for the archive (required)")
	_ = saveCmd.MarkFlagRequired("label")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved archives",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := os.Getwd()
			archives, err := envfile.ListArchives(base)
			if err != nil {
				return err
			}
			if len(archives) == 0 {
				fmt.Println("No archives found.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tLABEL\tCREATED\tKEYS")
			for _, a := range archives {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
					a.ID, a.Label, a.CreatedAt.Format("2006-01-02 15:04:05"), len(a.Entries))
			}
			return w.Flush()
		},
	}

	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore an archived env file by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if archiveID == "" {
				return fmt.Errorf("--id is required")
			}
			base, _ := os.Getwd()
			archive, err := envfile.LoadArchive(base, archiveID)
			if err != nil {
				return err
			}
			if err := envfile.Write(envPath, archive.Entries); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			fmt.Printf("Restored archive %q (%s) → %s\n", archive.Label, archive.ID, envPath)
			return nil
		},
	}
	restoreCmd.Flags().StringVar(&archiveID, "id", "", "Archive ID to restore")

	archiveCmd.AddCommand(saveCmd, listCmd, restoreCmd)
	rootCmd.AddCommand(archiveCmd)
}

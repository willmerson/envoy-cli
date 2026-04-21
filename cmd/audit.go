package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var auditCmd = &cobra.Command{
	Use:   "audit <before> <after>",
	Short: "Show what changed between two .env files",
	Long:  `Compare two .env files and display an audit log of added, removed, and changed keys.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		beforeEntries, err := loadAuditFile(args[0])
		if err != nil {
			return fmt.Errorf("loading before file: %w", err)
		}

		afterEntries, err := loadAuditFile(args[1])
		if err != nil {
			return fmt.Errorf("loading after file: %w", err)
		}

		log := envfile.Audit(beforeEntries, afterEntries)
		out := envfile.FormatAuditLog(log)
		fmt.Print(out)

		if jsonOut, _ := cmd.Flags().GetBool("json"); jsonOut {
			return printAuditJSON(log)
		}

		return nil
	},
}

func loadAuditFile(path string) ([]envfile.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return envfile.Parse(f)
}

func printAuditJSON(log envfile.AuditLog) error {
	for _, entry := range log {
		fmt.Printf(`{"timestamp":%q,"action":%q,"key":%q,"old_value":%q,"new_value":%q}`+"\n",
			entry.Timestamp.String(), entry.Action, entry.Key, entry.OldValue, entry.NewValue)
	}
	return nil
}

func init() {
	auditCmd.Flags().Bool("json", false, "Output audit log as JSON lines")
	rootCmd.AddCommand(auditCmd)
}

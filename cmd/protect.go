package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

func init() {
	var keys []string
	var prefixes []string
	var dryRun bool

	protectCmd := &cobra.Command{
		Use:   "protect",
		Short: "Mark keys as protected with sentinel comments",
		Long: `Inserts a #PROTECTED comment marker above specified keys to signal
that they should not be modified. Use --key and/or --prefix to select targets.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(keys) == 0 && len(prefixes) == 0 {
				return fmt.Errorf("at least one --key or --prefix must be provided")
			}

			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			// Flatten comma-separated keys
			var flatKeys []string
			for _, k := range keys {
				for _, part := range strings.Split(k, ",") {
					if t := strings.TrimSpace(part); t != "" {
						flatKeys = append(flatKeys, t)
					}
				}
			}

			updated, result, err := envfile.Protect(entries, envfile.ProtectOptions{
				Keys:     flatKeys,
				Prefixes: prefixes,
				DryRun:   dryRun,
			})
			if err != nil {
				return err
			}

			fmt.Print(envfile.FormatProtectResult(result))

			if dryRun {
				fmt.Fprintln(os.Stderr, "Dry-run: no changes written.")
				return nil
			}

			if err := envfile.Write(envPath, updated); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			return nil
		},
	}

	protectCmd.Flags().StringArrayVarP(&keys, "key", "k", nil, "Key name(s) to protect (repeatable, comma-separated)")
	protectCmd.Flags().StringArrayVarP(&prefixes, "prefix", "p", nil, "Key prefix(es) to protect (repeatable)")
	protectCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")

	rootCmd.AddCommand(protectCmd)
}

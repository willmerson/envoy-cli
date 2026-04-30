package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	scopePrefix      string
	scopeStrip       bool
	scopeCaseSensitive bool
	scopeOutput      string
)

func init() {
	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Filter entries by key prefix and optionally strip the prefix",
		Example: `  envoy scope --prefix APP_ --strip
  envoy scope --prefix DB_ --output db.env`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			result := envfile.Scope(entries, envfile.ScopeOptions{
				Prefix:        scopePrefix,
				StripPrefix:   scopeStrip,
				CaseSensitive: scopeCaseSensitive,
			})

			dest := envPath
			if scopeOutput != "" {
				dest = scopeOutput
			}

			if scopeOutput != "" {
				if err := envfile.Write(dest, result.Entries); err != nil {
					return fmt.Errorf("write: %w", err)
				}
			} else {
				if err := envfile.WriteTo(os.Stdout, result.Entries); err != nil {
					return fmt.Errorf("write: %w", err)
				}
			}

			fmt.Fprintln(os.Stderr, envfile.FormatScopeSummary(result, scopePrefix))
			_ = dest
			return nil
		},
	}

	scopeCmd.Flags().StringVar(&scopePrefix, "prefix", "", "Key prefix to scope entries by")
	scopeCmd.Flags().BoolVar(&scopeStrip, "strip", false, "Strip the prefix from matched keys")
	scopeCmd.Flags().BoolVar(&scopeCaseSensitive, "case-sensitive", false, "Use case-sensitive prefix matching")
	scopeCmd.Flags().StringVar(&scopeOutput, "output", "", "Write scoped entries to this file instead of stdout")

	rootCmd.AddCommand(scopeCmd)
}

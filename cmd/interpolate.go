package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/envfile"
)

func init() {
	var strict bool
	var outputPath string

	interpolateCmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Expand variable references within .env values",
		Long: `Resolve $VAR and ${VAR} references inside .env values using
other keys defined in the same file. Optionally fail on unresolved references
with --strict.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			res, err := envfile.Interpolate(entries, strict)
			if err != nil {
				return err
			}

			dest := envPath
			if outputPath != "" {
				dest = outputPath
			}

			if err := envfile.Write(dest, res.Entries); err != nil {
				return fmt.Errorf("write: %w", err)
			}

			fmt.Fprintf(os.Stderr, "interpolate: %d reference(s) expanded", res.Expanded)
			if len(res.Unresolved) > 0 {
				fmt.Fprintf(os.Stderr, ", unresolved: %s", strings.Join(res.Unresolved, ", "))
			}
			fmt.Fprintln(os.Stderr)
			return nil
		},
	}

	interpolateCmd.Flags().BoolVar(&strict, "strict", false, "fail if any references cannot be resolved")
	interpolateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "write result to this file instead of the source file")

	rootCmd.AddCommand(interpolateCmd)
}

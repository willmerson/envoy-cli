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
	var patterns []string
	var placeholder string
	var showLast int
	var outputPath string

	maskCmd := &cobra.Command{
		Use:   "mask",
		Short: "Mask sensitive values in a .env file",
		Long: `Mask replaces sensitive values with a placeholder.
Specify keys by exact name (--key) or substring pattern (--pattern).
Use --show-last N to reveal the last N characters of a masked value.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envfile.MaskOptions{
				Keys:        keys,
				Patterns:    patterns,
				Placeholder: placeholder,
				ShowLast:    showLast,
			}

			result := envfile.Mask(entries, opts)

			dest := envPath
			if outputPath != "" {
				dest = outputPath
			}

			if err := envfile.Write(dest, result.Entries); err != nil {
				return fmt.Errorf("write: %w", err)
			}

			fmt.Fprintln(os.Stdout, envfile.FormatMaskSummary(result))
			return nil
		},
	}

	maskCmd.Flags().StringArrayVar(&keys, "key", nil, "exact key name(s) to mask (repeatable)")
	maskCmd.Flags().StringArrayVar(&patterns, "pattern", nil, "substring pattern(s) to match key names (repeatable)")
	maskCmd.Flags().StringVar(&placeholder, "placeholder", "****", "replacement placeholder for masked values")
	maskCmd.Flags().IntVar(&showLast, "show-last", 0, "reveal last N characters of masked value")
	maskCmd.Flags().StringVarP(&outputPath, "output", "o", "", "write result to this file instead of the source")

	_ = maskCmd.MarkFlagRequired

	rootCmd.AddCommand(maskCmd)

	_ = strings.ToUpper // satisfy import
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	redactKeys       []string
	redactPatterns   []string
	redactPlaceholder string
	redactInPlace    bool
)

var redactCmd = &cobra.Command{
	Use:   "redact [file]",
	Short: "Redact sensitive values in a .env file",
	Long: `Replace sensitive values with a placeholder.
Match keys by exact name (--key) or regex pattern (--pattern).
Outputs to stdout by default; use --in-place to overwrite the file.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := envPath
		if len(args) > 0 {
			path = args[0]
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}

		entries, err := envfile.Parse(string(data))
		if err != nil {
			return fmt.Errorf("parsing file: %w", err)
		}

		result := envfile.Redact(entries, envfile.RedactOptions{
			Keys:        redactKeys,
			Patterns:    redactPatterns,
			Placeholder: redactPlaceholder,
		})

		if redactInPlace {
			if err := envfile.Write(path, result.Entries); err != nil {
				return fmt.Errorf("writing file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Redacted %d key(s) in %s\n", result.Redacted, path)
			return nil
		}

		var sb strings.Builder
		if err := envfile.WriteTo(&sb, result.Entries); err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), sb.String())
		return nil
	},
}

func init() {
	redactCmd.Flags().StringArrayVarP(&redactKeys, "key", "k", nil, "Key name(s) to redact (case-insensitive)")
	redactCmd.Flags().StringArrayVarP(&redactPatterns, "pattern", "p", nil, "Regex pattern(s) matched against key names")
	redactCmd.Flags().StringVar(&redactPlaceholder, "placeholder", "***", "Replacement value for redacted entries")
	redactCmd.Flags().BoolVarP(&redactInPlace, "in-place", "i", false, "Overwrite the source file")
	rootCmd.AddCommand(redactCmd)
}

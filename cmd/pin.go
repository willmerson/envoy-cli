package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

func init() {
	var (
		keys      []string
		prefix    string
		overwrite bool
		dryRun    bool
	)

	pinCmd := &cobra.Command{
		Use:   "pin [file] KEY=VALUE...",
		Short: "Pin one or more keys to specific values in an .env file",
		Long: `Pin locks the values of selected keys to the provided values.
Keys already holding the target value are skipped unless --overwrite is set.

Example:
  envoy pin .env APP_VERSION=2.3.1 DB_HOST=prod-db
  envoy pin .env --keys APP_VERSION,DB_HOST --prefix DB_ APP_VERSION=2.3.1`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			pairs := args[1:]

			if len(pairs) == 0 {
				return fmt.Errorf("at least one KEY=VALUE pair is required")
			}

			pinMap := make(map[string]string, len(pairs))
			for _, p := range pairs {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid pair %q: expected KEY=VALUE", p)
				}
				pinMap[parts[0]] = parts[1]
			}

			entries, err := envfile.Parse(filePath)
			if err != nil {
				return fmt.Errorf("parse %s: %w", filePath, err)
			}

			opts := envfile.PinOptions{
				Keys:      keys,
				Prefix:    prefix,
				Overwrite: overwrite,
			}

			updated, results := envfile.Pin(entries, pinMap, opts)
			fmt.Print(envfile.FormatPinSummary(results))

			if dryRun {
				fmt.Fprintln(os.Stderr, "dry-run: no changes written")
				return nil
			}

			if err := envfile.Write(filePath, updated); err != nil {
				return fmt.Errorf("write %s: %w", filePath, err)
			}
			return nil
		},
	}

	pinCmd.Flags().StringSliceVar(&keys, "keys", nil, "Comma-separated list of keys to pin (default: all in pinMap)")
	pinCmd.Flags().StringVar(&prefix, "prefix", "", "Only pin keys with this prefix")
	pinCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite even if value is already set to target")
	pinCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")

	rootCmd.AddCommand(pinCmd)
}

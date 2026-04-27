package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

func init() {
	var keyPairs []string
	var oldPrefix, newPrefix string
	var failOnMissing bool
	var dryRun bool

	rotateCmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate (rename) keys in a .env file",
		Long: `Rotate renames keys in a .env file either by an explicit key map
or by replacing a common prefix across matching keys.`,
		Example: `  envoy rotate --map OLD_KEY=NEW_KEY --map APP_SECRET=SVC_SECRET
  envoy rotate --old-prefix APP_ --new-prefix SVC_`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envfile.RotateOptions{
				OldPrefix:     oldPrefix,
				NewPrefix:     newPrefix,
				FailOnMissing: failOnMissing,
			}

			if len(keyPairs) > 0 {
				opts.KeyMap = make(map[string]string, len(keyPairs))
				for _, pair := range keyPairs {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
						return fmt.Errorf("invalid key pair %q: expected OLD=NEW", pair)
					}
					opts.KeyMap[parts[0]] = parts[1]
				}
			}

			out, result, err := envfile.Rotate(entries, opts)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stderr, envfile.FormatRotateResult(result))

			if dryRun {
				return nil
			}

			if err := envfile.Write(envPath, out); err != nil {
				return fmt.Errorf("write: %w", err)
			}
			return nil
		},
	}

	rotateCmd.Flags().StringArrayVar(&keyPairs, "map", nil, "Key pair to rotate in OLD=NEW format (repeatable)")
	rotateCmd.Flags().StringVar(&oldPrefix, "old-prefix", "", "Prefix to replace")
	rotateCmd.Flags().StringVar(&newPrefix, "new-prefix", "", "Replacement prefix")
	rotateCmd.Flags().BoolVar(&failOnMissing, "fail-on-missing", false, "Error if a mapped key is not found")
	rotateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print result without writing changes")

	rootCmd.AddCommand(rotateCmd)
}

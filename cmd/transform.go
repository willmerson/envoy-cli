package cmd

import (
	"fmt"
	"strings"

	"github.com/saurabh/envoy-cli/internal/envfile"
	"github.com/spf13/cobra"
)

func init() {
	var target string
	var ops []string

	transformCmd := &cobra.Command{
		Use:   "transform",
		Short: "Apply transformations to keys or values in a .env file",
		Long: `Transform keys or values using operations: uppercase, lowercase, trimspace, quoteall, unquoteall.
Multiple operations are applied in order.`,
		Example: `  envoy transform --ops uppercase --target keys
  envoy transform --ops trimspace,lowercase --target values`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("failed to parse env file: %w", err)
			}

			applyToKeys := strings.EqualFold(target, "keys")

			var opts []envfile.TransformOption
			for _, op := range ops {
				switch strings.ToLower(strings.TrimSpace(op)) {
				case "uppercase":
					opts = append(opts, envfile.TransformUppercase)
				case "lowercase":
					opts = append(opts, envfile.TransformLowercase)
				case "trimspace":
					opts = append(opts, envfile.TransformTrimSpace)
				case "quoteall":
					opts = append(opts, envfile.TransformQuoteAll)
				case "unquoteall":
					opts = append(opts, envfile.TransformUnquoteAll)
				default:
					return fmt.Errorf("unknown operation: %s", op)
				}
			}

			if len(opts) == 0 {
				return fmt.Errorf("at least one operation must be specified via --ops")
			}

			result, summary := envfile.Transform(entries, applyToKeys, opts...)

			if err := envfile.Write(envPath, result); err != nil {
				return fmt.Errorf("failed to write env file: %w", err)
			}

			fmt.Printf("Transformed %d/%d entries (target: %s, ops: %s)\n",
				summary.Modified, summary.Total, target, strings.Join(ops, ","))
			return nil
		},
	}

	transformCmd.Flags().StringVar(&target, "target", "values", "What to transform: keys or values")
	transformCmd.Flags().StringSliceVar(&ops, "ops", nil, "Comma-separated list of operations to apply")
	_ = transformCmd.MarkFlagRequired("ops")

	rootCmd.AddCommand(transformCmd)
}

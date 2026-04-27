package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/envoy-cli/internal/envfile"
)

func init() {
	var (
		genType      string
		defaultValue string
		randomLen    int
		overwrite    bool
		outPath      string
	)

	cmd := &cobra.Command{
		Use:   "generate [flags] KEY [KEY...]",
		Short: "Generate placeholder or random values for env keys",
		Long: `Add new keys with generated values to an .env file.

Types:
  literal  Use --default as the value (default: empty string)
  random   Cryptographically random hex string (see --length)
  uuid     Random UUID v4`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("env")
			if outPath == "" {
				outPath = path
			}

			var base []envfile.Entry
			if parsed, err := envfile.ParseFile(path); err == nil {
				base = parsed
			}

			opts := envfile.GenerateOptions{
				Keys:         args,
				Type:         genType,
				DefaultValue: defaultValue,
				RandomLength: randomLen,
				Overwrite:    overwrite,
			}

			res, err := envfile.Generate(base, opts)
			if err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			if err := envfile.Write(outPath, res.Entries); err != nil {
				return fmt.Errorf("write: %w", err)
			}

			if len(res.Added) > 0 {
				fmt.Printf("generated: %s\n", strings.Join(res.Added, ", "))
			}
			if len(res.Skipped) > 0 {
				fmt.Printf("skipped (already exist): %s\n", strings.Join(res.Skipped, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&genType, "type", "literal", "Value type: literal, random, uuid")
	cmd.Flags().StringVar(&defaultValue, "default", "", "Default value for literal type")
	cmd.Flags().IntVar(&randomLen, "length", 16, "Byte length for random type")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys")
	cmd.Flags().StringVar(&outPath, "out", "", "Output file (defaults to --env path)")

	rootCmd.AddCommand(cmd)
}

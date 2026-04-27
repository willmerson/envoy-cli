package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

func init() {
	var outputDir string
	var separator string
	var lowercase bool
	var dryRun bool

	splitCmd := &cobra.Command{
		Use:   "split",
		Short: "Split a .env file into multiple files grouped by key prefix",
		Example: `  envoy split --output ./envs
  envoy split --sep . --lowercase --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			opts := envfile.SplitOptions{
				ByPrefix:  true,
				Separator: separator,
				Lowercase: lowercase,
			}
			result := envfile.Split(entries, opts)

			fmt.Print(envfile.FormatSplitSummary(result))

			if dryRun {
				fmt.Println("(dry-run: no files written)")
				return nil
			}

			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return fmt.Errorf("cannot create output directory: %w", err)
			}

			for label, grpEntries := range result.Groups {
				fileName := label + ".env"
				outPath := filepath.Join(outputDir, fileName)
				if err := envfile.Write(grpEntries, outPath); err != nil {
					return fmt.Errorf("write %s: %w", outPath, err)
				}
				fmt.Printf("  wrote %s\n", outPath)
			}

			if len(result.Ungrouped) > 0 {
				outPath := filepath.Join(outputDir, "other.env")
				if err := envfile.Write(result.Ungrouped, outPath); err != nil {
					return fmt.Errorf("write %s: %w", outPath, err)
				}
				fmt.Printf("  wrote %s\n", outPath)
			}

			return nil
		},
	}

	splitCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Directory to write split .env files into")
	splitCmd.Flags().StringVar(&separator, "sep", "_", "Key separator used to detect prefixes")
	splitCmd.Flags().BoolVar(&lowercase, "lowercase", false, "Lowercase output filenames")
	splitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview groups without writing files")

	rootCmd.AddCommand(splitCmd)
}

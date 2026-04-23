package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var (
	promoteOverwrite    bool
	promotePrefixFilter string
	promoteStripPrefix  bool
	promoteOutput       string
)

var promoteCmd = &cobra.Command{
	Use:   "promote <source-env> <target-env>",
	Short: "Promote entries from one env file into another",
	Long: `Copies entries from a source .env file into a target .env file.
By default, existing keys in the target are not overwritten.
Use --overwrite to replace them.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		dstPath := args[1]

		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("reading source: %w", err)
		}
		srcEntries, err := envfile.Parse(string(srcData))
		if err != nil {
			return fmt.Errorf("parsing source: %w", err)
		}

		var dstEntries []envfile.Entry
		dstData, err := os.ReadFile(dstPath)
		if err == nil {
			dstEntries, err = envfile.Parse(string(dstData))
			if err != nil {
				return fmt.Errorf("parsing target: %w", err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("reading target: %w", err)
		}

		opts := envfile.PromoteOptions{
			Overwrite:    promoteOverwrite,
			PrefixFilter: promotePrefixFilter,
			StripPrefix:  promoteStripPrefix,
		}

		result, res, err := envfile.Promote(srcEntries, dstEntries, opts)
		if err != nil {
			return err
		}

		outPath := dstPath
		if promoteOutput != "" {
			outPath = promoteOutput
		}

		if err := envfile.Write(outPath, result); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), envfile.FormatPromoteResult(res))
		return nil
	},
}

func init() {
	promoteCmd.Flags().BoolVar(&promoteOverwrite, "overwrite", false, "Overwrite existing keys in target")
	promoteCmd.Flags().StringVar(&promotePrefixFilter, "prefix", "", "Only promote keys with this prefix")
	promoteCmd.Flags().BoolVar(&promoteStripPrefix, "strip-prefix", false, "Strip prefix from keys when writing to target")
	promoteCmd.Flags().StringVarP(&promoteOutput, "output", "o", "", "Write result to this file instead of target")
	rootCmd.AddCommand(promoteCmd)
}

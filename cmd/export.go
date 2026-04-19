package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/spf13/cobra"
)

var exportFormat string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export .env file in a specified format (dotenv, json, export)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := envPath
		if len(args) > 0 {
			path = args[0]
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening env file: %w", err)
		}
		defer f.Close()

		entries, err := envfile.Parse(f)
		if err != nil {
			return fmt.Errorf("parsing env file: %w", err)
		}

		out, err := envfile.Export(entries, envfile.ExportFormat(exportFormat))
		if err != nil {
			return err
		}

		fmt.Print(out)
		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, json, export")
	rootCmd.AddCommand(exportCmd)
}

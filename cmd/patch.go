package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/umutdz/envoy-cli/internal/envfile"
)

var patchFile string

type patchSpec struct {
	Ops []envfile.PatchOp `json:"ops"`
}

func init() {
	patchCmd := &cobra.Command{
		Use:   "patch",
		Short: "Apply a set of patch operations (set/delete/rename) to a .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := os.ReadFile(patchFile)
			if err != nil {
				return fmt.Errorf("reading patch file: %w", err)
			}

			var spec patchSpec
			if err := json.Unmarshal(raw, &spec); err != nil {
				return fmt.Errorf("parsing patch file: %w", err)
			}

			entries, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parsing env file: %w", err)
			}

			result, patchResult := envfile.Patch(entries, spec.Ops)

			for _, a := range patchResult.Applied {
				fmt.Fprintf(cmd.OutOrStdout(), "applied: %s\n", a)
			}
			for _, s := range patchResult.Skipped {
				fmt.Fprintf(cmd.OutOrStdout(), "skipped: %s\n", s)
			}
			for _, e := range patchResult.Errors {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", e)
			}

			if err := envfile.Write(envPath, result); err != nil {
				return fmt.Errorf("writing env file: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "patch applied to %s\n", envPath)
			return nil
		},
	}

	patchCmd.Flags().StringVarP(&patchFile, "patch-file", "p", "", "Path to JSON patch file (required)")
	_ = patchCmd.MarkFlagRequired("patch-file")

	rootCmd.AddCommand(patchCmd)
}

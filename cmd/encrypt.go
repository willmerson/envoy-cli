package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var encryptPassphrase string

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt all values in a .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
		for i, e := range entries {
			if e.Value == "" {
				continue
			}
			enc, err := envfile.Encrypt(e.Value, encryptPassphrase)
			if err != nil {
				return fmt.Errorf("encrypt %s: %w", e.Key, err)
			}
			entries[i].Value = enc
		}
		if err := envfile.Write(envPath, entries); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Encrypted %d entries in %s\n", len(entries), envPath)
		return nil
	},
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt all values in a .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := envfile.Parse(envPath)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
		for i, e := range entries {
			if e.Value == "" {
				continue
			}
			dec, err := envfile.Decrypt(e.Value, encryptPassphrase)
			if err != nil {
				return fmt.Errorf("decrypt %s: %w", e.Key, err)
			}
			entries[i].Value = dec
		}
		if err := envfile.Write(envPath, entries); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Decrypted %d entries in %s\n", len(entries), envPath)
		return nil
	},
}

func init() {
	encryptCmd.Flags().StringVarP(&encryptPassphrase, "passphrase", "p", "", "Passphrase for encryption (required)")
	_ = encryptCmd.MarkFlagRequired("passphrase")
	decryptCmd.Flags().StringVarP(&encryptPassphrase, "passphrase", "p", "", "Passphrase for decryption (required)")
	_ = decryptCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/envfile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage environment profiles",
}

var profileSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current .env as a named profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		f, err := os.Open(envPath)
		if err != nil {
			return fmt.Errorf("opening %s: %w", envPath, err)
		}
		defer f.Close()

		pairs, err := envfile.Parse(f)
		if err != nil {
			return err
		}

		p := &envfile.Profile{Name: name}
		for k, v := range pairs {
			p.Entries = append(p.Entries, envfile.Entry{Key: k, Value: v})
		}

		cwd, _ := os.Getwd()
		if err := envfile.SaveProfile(cwd, p); err != nil {
			return err
		}
		fmt.Printf("Profile %q saved.\n", name)
		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		names, err := envfile.ListProfiles(cwd)
		if err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Println("No profiles found.")
			return nil
		}
		for _, n := range names {
			fmt.Println(n)
		}
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileSaveCmd)
	profileCmd.AddCommand(profileListCmd)
	rootCmd.AddCommand(profileCmd)
}

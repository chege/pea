package cmd

import (
	"fmt"
	"pea/internal/app"

	"github.com/spf13/cobra"
)

func addRemoteCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "remote <url>",
		Short: "Configure a remote git repository for synchronization",
		Long:  "Configure a remote git repository (GitHub, GitLab, Gist, etc.) to sync your prompts.\nThis updates your ~/.pea/config.toml and the git configuration in your store.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			// 1. Update config file
			if err := app.UpdateRemoteURL(url); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}

			// 2. Ensure store is loaded and get path
			// Note: EnsureStore re-reads the config we just wrote.
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}

			// 3. Configure git remote (force update)
			if err := app.SetGitRemote(store, url); err != nil {
				return err
			}

			fmt.Printf("Remote configured: %s\n", url)
			fmt.Println("Future changes will be automatically pushed.")
			return nil
		},
	}
	root.AddCommand(cmd)
}

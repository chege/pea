package cmd

import (
	"fmt"
	"pea/internal/app"

	"github.com/spf13/cobra"
)

func addSyncCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Manually sync with remote git repository",
		Long:  "Perform a manual git pull --rebase and git push to synchronize with the configured remote.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}

			fmt.Println("Syncing with remote...")
			if err := app.Sync(store, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
				return err
			}
			fmt.Println("Sync complete.")
			return nil
		},
	}
	root.AddCommand(cmd)
}

package cmd

import (
	"fmt"
	"pea/internal/app"

	"github.com/spf13/cobra"
)

func addListCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list stored entries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}
			entries, err := app.ListEntries(store)
			if err != nil {
				return err
			}
			for _, e := range entries {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), e); err != nil {
					return err
				}
			}
			return nil
		},
	}
	root.AddCommand(cmd)
}

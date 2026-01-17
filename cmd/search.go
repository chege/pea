package cmd

import (
	"fmt"
	"pea/internal/app"
	"strings"

	"github.com/spf13/cobra"
)

func addSearchCommand(root *cobra.Command) {
	var tags []string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "search entries by name substring or tags",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}
			query := ""
			if len(args) == 1 {
				query = strings.ToLower(args[0])
			}

			entries, err := app.CollectEntriesWithTags(store)
			if err != nil {
				return err
			}

			var out []string
			for _, e := range entries {
				match := false
				if query == "" {
					match = true
				} else {
					if strings.Contains(strings.ToLower(e.Name), query) {
						match = true
					} else if strings.Contains(strings.ToLower(e.Content), query) {
						match = true
					}
				}

				if match && len(tags) > 0 && !app.HasAllTags(e.Tags, tags) {
					match = false
				}

				if match {
					out = append(out, e.Name)
				}
			}

			for _, name := range out {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "filter by tag (repeatable)")
	root.AddCommand(cmd)
}

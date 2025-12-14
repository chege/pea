package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func addRemoveCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:               "rm <name>",
		Short:             "delete an entry",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			name, err := normalizeName(args[0])
			if err != nil {
				return err
			}
			path, ext, err := existingEntryPath(store, name)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("delete failed: not found: %s", name)
				}
				return err
			}
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}
			// git rm + commit (best-effort)
			commitMsg := "chore: remove " + name + ext
			gitRmAndCommit(store, []string{name + ext}, commitMsg, cmd.ErrOrStderr())
			_, err = fmt.Fprintln(cmd.OutOrStdout(), name)
			return err
		},
	}
	root.AddCommand(cmd)
}

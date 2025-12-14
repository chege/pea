package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func addMoveCommand(root *cobra.Command) {
	var choreRename bool

	cmd := &cobra.Command{
		Use:   "mv <old> <new>",
		Short: "rename an entry",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) >= 1 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return completeNames(cmd, args, toComplete)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			oldName, err := normalizeName(args[0])
			if err != nil {
				return err
			}
			newName, err := normalizeName(args[1])
			if err != nil {
				return err
			}
			oldPath, ext, err := existingEntryPath(store, oldName)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("rename failed: not found: %s", oldName)
				}
				return fmt.Errorf("rename failed: %w", err)
			}
			newPath := defaultEntryPath(store, newName)
			if ext == legacyExt {
				newPath = legacyEntryPath(store, newName)
			}
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}
			// git add new and commit (best-effort)
			commitMsg := fmt.Sprintf("refactor: rename %s%s to %s%s", oldName, ext, newName, ext)
			if choreRename {
				commitMsg = fmt.Sprintf("chore: rename %s%s to %s%s", oldName, ext, newName, ext)
			}
			gitAddAndCommit(store, []string{oldName + ext, newName + ext}, commitMsg, cmd.ErrOrStderr())
			_, err = fmt.Fprintln(cmd.OutOrStdout(), newName)
			return err
		},
	}
	cmd.Flags().BoolVar(&choreRename, "chore", false, "mark rename as organizational")
	root.AddCommand(cmd)
}

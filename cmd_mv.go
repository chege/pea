package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func addMoveCommand(root *cobra.Command) {
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
			oldName := toSnake(args[0])
			newName := toSnake(args[1])
			oldPath, ext, err := existingEntryPath(store, oldName)
			if err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}
			newPath := defaultEntryPath(store, newName)
			if ext == legacyExt {
				newPath = legacyEntryPath(store, newName)
			}
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}
			// git add new and commit
			commitMsg := fmt.Sprintf("refactor: rename %s%s to %s%s", oldName, ext, newName, ext)
			_ = exec.Command("bash", "-c", "cd '"+store+"' && git add '"+newName+ext+"' && git commit -m '"+commitMsg+"'").Run()
			_, err = fmt.Fprintln(cmd.OutOrStdout(), newName)
			return err
		},
	}
	root.AddCommand(cmd)
}

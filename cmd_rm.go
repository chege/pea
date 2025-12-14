package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func addRemoveCommand(root *cobra.Command) {
	var confirm bool
	var dryRun bool
	var undo bool

	cmd := &cobra.Command{
		Use:               "rm <name>",
		Short:             "delete an entry",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			if undo {
				store, err := ensureStore()
				if err != nil {
					return err
				}
				name, err := normalizeName(args[0])
				if err != nil {
					return err
				}
				_, ext, err := existingEntryPath(store, name)
				if err != nil && !os.IsNotExist(err) {
					return err
				}
				if ext == "" {
					ext = defaultExt
				}
				if err := revertLastCommitForPath(store, name+ext, cmd.ErrOrStderr()); err != nil {
					return fmt.Errorf("undo failed: %w", err)
				}
				_, err = fmt.Fprintln(cmd.OutOrStdout(), name)
				return err
			}
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
			if dryRun {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "dry-run: would delete %s%s\n", name, ext)
				return err
			}
			if confirm {
				fmt.Fprintf(cmd.OutOrStdout(), "Delete %s%s? [y/N]: ", name, ext)
				reader := bufio.NewReader(cmd.InOrStdin())
				ans, _ := reader.ReadString('\n')
				ans = strings.ToLower(strings.TrimSpace(ans))
				if ans != "y" && ans != "yes" {
					return fmt.Errorf("delete aborted")
				}
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
	cmd.Flags().BoolVar(&confirm, "confirm", false, "prompt before deleting")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would happen without deleting")
	cmd.Flags().BoolVar(&undo, "undo", false, "undo last delete for this entry via git revert")
	root.AddCommand(cmd)
}

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"pea/internal/app"
	"strings"

	"github.com/spf13/cobra"
)

func addMoveCommand(root *cobra.Command) {
	var choreRename bool
	var confirm bool
	var dryRun bool

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
			store, err := app.EnsureStore()
			if err != nil {
				return err
			}
			oldName, err := app.NormalizeName(args[0])
			if err != nil {
				return err
			}
			newName, err := app.NormalizeName(args[1])
			if err != nil {
				return err
			}
			oldPath, ext, err := app.ExistingEntryPath(store, oldName)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("rename failed: not found: %s", oldName)
				}
				return fmt.Errorf("rename failed: %w", err)
			}
			newPath := app.DefaultEntryPath(store, newName)
			if ext == app.LegacyExt {
				newPath = app.LegacyEntryPath(store, newName)
			}
			if dryRun {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "dry-run: would rename %s%s to %s%s\n", oldName, ext, newName, ext)
				return err
			}
			if confirm {
				fmt.Fprintf(cmd.OutOrStdout(), "Rename %s%s to %s%s? [y/N]: ", oldName, ext, newName, ext)
				reader := bufio.NewReader(cmd.InOrStdin())
				ans, _ := reader.ReadString('\n')
				ans = strings.ToLower(strings.TrimSpace(ans))
				if ans != "y" && ans != "yes" {
					return fmt.Errorf("rename aborted")
				}
			}
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}
			// git add new and commit (best-effort)
			commitMsg := fmt.Sprintf("refactor: rename %s%s to %s%s", oldName, ext, newName, ext)
			if choreRename {
				commitMsg = fmt.Sprintf("chore: rename %s%s to %s%s", oldName, ext, newName, ext)
			}
			app.GitAddAndCommit(store, []string{oldName + ext, newName + ext}, commitMsg, cmd.ErrOrStderr())
			_, err = fmt.Fprintln(cmd.OutOrStdout(), newName)
			return err
		},
	}
	cmd.Flags().BoolVar(&choreRename, "chore", false, "mark rename as organizational")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "prompt before renaming")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would happen without renaming")
	root.AddCommand(cmd)
}

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"pea/internal/app"

	"github.com/spf13/cobra"
)

func addEditCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "edit a snippet in $EDITOR",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(cmd, args[0])
		},
	}
	root.AddCommand(cmd)
}

func runEdit(cmd *cobra.Command, nameRaw string) error {
	store, err := app.EnsureStore()
	if err != nil {
		return err
	}

	name, err := app.NormalizeName(nameRaw)
	if err != nil {
		return err
	}

	path, ext, err := app.ExistingEntryPath(store, name)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("entry not found: %s", name)
		}
		return err
	}

	// Read original content
	originalContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read entry: %w", err)
	}

	// Open editor
	if err := openEditor(cmd, path); err != nil {
		return err
	}

	// Read new content
	newContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read entry after edit: %w", err)
	}

	// Check for changes
	if bytes.Equal(originalContent, newContent) {
		fmt.Fprintln(cmd.OutOrStdout(), "no changes")
		return nil
	}

	if len(bytes.TrimSpace(newContent)) == 0 {
		return fmt.Errorf("edit aborted: empty content")
	}

	// Commit
	commitMsg := "feat: edit " + name + ext
	app.GitAddAndCommit(store, []string{name + ext}, commitMsg, cmd.ErrOrStderr())

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", name)
	return nil
}

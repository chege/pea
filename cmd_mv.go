package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func addMoveCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "mv <old> <new>",
		Short: "rename an entry",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil { return err }
			oldName := toSnake(args[0])
			newName := toSnake(args[1])
			oldPath := filepath.Join(store, oldName+".txt")
			newPath := filepath.Join(store, newName+".txt")
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), newName)
			return nil
		},
	}
	root.AddCommand(cmd)
}

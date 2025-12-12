package main

import (
	"fmt"
	"os"
	"os/exec"
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
			if err != nil {
				return err
			}
			oldName := toSnake(args[0])
			newName := toSnake(args[1])
			oldPath := filepath.Join(store, oldName+".txt")
			newPath := filepath.Join(store, newName+".txt")
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}
			// git add new and commit (best-effort)
			_ = exec.Command("bash", "-c", "cd '"+store+"' && git add '"+newName+".txt' && git commit -m 'mv "+oldName+" -> "+newName+"'").Run()
			fmt.Fprintln(cmd.OutOrStdout(), newName)
			return nil
		},
	}
	root.AddCommand(cmd)
}

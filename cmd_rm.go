package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func addRemoveCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "rm <name>",
		Short: "delete an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ensureStore()
			if err != nil {
				return err
			}
			name := toSnake(args[0])
			path := filepath.Join(store, name+".txt")
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}
			// git rm + commit
			_ = exec.Command("bash", "-c", "cd '"+store+"' && git rm -f '"+name+".txt' && git commit -m 'rm "+name+"'").Run()
			fmt.Fprintln(cmd.OutOrStdout(), name)
			return nil
		},
	}
	root.AddCommand(cmd)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func addCompletionCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|install]",
		Short: "output or install shell completion",
		Args:  cobra.MinimumNArgs(1),
		ValidArgs: []string{"bash", "zsh", "install"},
		RunE: func(cmd *cobra.Command, args []string) error {
			op := args[0]
			switch op {
			case "bash":
				return root.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "install":
				home, err := os.UserHomeDir()
				if err != nil { return err }
				base := filepath.Join(home, ".pea")
				_ = os.MkdirAll(base, 0o755)
				bashPath := filepath.Join(base, "p.bash")
				zshPath := filepath.Join(base, "_p")
				{
					f, err := os.Create(bashPath)
					if err != nil { return err }
					defer f.Close()
					if err := root.GenBashCompletion(f); err != nil { return err }
				}
				{
					f, err := os.Create(zshPath)
					if err != nil { return err }
					defer f.Close()
					if err := root.GenZshCompletion(f); err != nil { return err }
				}
				fmt.Fprintf(cmd.OutOrStdout(), "installed completion: bash=%s zsh=%s\n", bashPath, zshPath)
				return nil
			default:
				return cmd.Help()
			}
		},
	}
	root.AddCommand(cmd)
}

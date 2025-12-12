package main

import (
	"github.com/spf13/cobra"
)

func addCompletionCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh]",
		Short: "output shell completion script",
		Args:  cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			switch shell {
			case "bash":
				return root.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			default:
				return cmd.Help()
			}
		},
	}
	root.AddCommand(cmd)
}

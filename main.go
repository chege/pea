package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "p",
		Short: "p: fast local prompt storage & retrieval",
		Long:  "p is a fast, local CLI to store short text under names and retrieve it instantly.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default behavior: show help
			return cmd.Help()
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = false

	cmd.Version = version
	cmd.SetVersionTemplate("p version {{.Version}}\n")

	addListCommand(cmd)
	return cmd
}

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

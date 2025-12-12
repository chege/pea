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
		Use:   "p [name]",
		Short: "p: fast local prompt storage & retrieval",
		Long:  "p is a fast, local CLI to store short text under names and retrieve it instantly.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				store, err := ensureStore()
				if err != nil { return err }
				b, err := readEntry(store, args[0])
				if err != nil { return err }
				_, err = cmd.OutOrStdout().Write(b)
				return err
			}
			return cmd.Help()
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = false

	cmd.Version = version
	cmd.SetVersionTemplate("p version {{.Version}}\n")

	addListCommand(cmd)
	addAddCommand(cmd)
	addRemoveCommand(cmd)
	return cmd
}

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
